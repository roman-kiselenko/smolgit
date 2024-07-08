package ssh

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"smolgit/internal/db"

	"golang.org/x/crypto/ssh"
)

// Code below partially based on https://github.com/gogs/gogs/blob/main/internal/ssh/ssh.go

func Listen(logger *slog.Logger, db *db.Database, opts conf.SSHOpts, appDataPath string) {
	config := &ssh.ServerConfig{
		Config: ssh.Config{},
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			pkey, err := db.FindPubKey(strings.TrimSpace(string(ssh.MarshalAuthorizedKey(key))))
			if err != nil {
				logger.Error("cant find pubkey", "err", err)
				return nil, err
			}
			return &ssh.Permissions{Extensions: map[string]string{"key-id": pkey.ID}}, nil
		},
	}

	keys, err := setupHostKeys(logger, appDataPath, opts.ServerAlgorithms)
	if err != nil {
		logger.Error("Failed to setup host keys", "err", err)
		os.Exit(1)

	}
	for _, key := range keys {
		config.AddHostKey(key)
	}

	go listen(config, opts.ListenHost, opts.ListenPort)
}

func listen(logger *slog.Logger, config *ssh.ServerConfig, host string, port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		logger.Error("Failed to start SSH server", "err", err)
		os.Exit(1)
	}
	for {
		// Once a ServerConfig has been configured, connections can be accepted.
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("Error accepting incoming connection", "err", err)
			continue
		}

		go func() {
			logger.Debug("Handshaking for %s", "addr", conn.RemoteAddr())
			sConn, chans, reqs, err := ssh.NewServerConn(conn, config)
			if err != nil {
				if err == io.EOF || errors.Is(err, syscall.ECONNRESET) {
					logger.Debug("Handshaking was terminated", "err", err)
				} else {
					logger.Debug("Error on handshaking", "err", err)
				}
				return
			}

			logger.Debug("Got connection", "from", sConn.RemoteAddr(), "version", sConn.ClientVersion())
			go ssh.DiscardRequests(reqs)
			go handleServerConn(logger, sConn.Permissions.Extensions["key-id"], chans)
		}()
	}
}

func handleServerConn(logger *slog.Logger, keyID string, chans <-chan ssh.NewChannel) {
	for newChan := range chans {
		if newChan.ChannelType() != "session" {
			_ = newChan.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		ch, reqs, err := newChan.Accept()
		if err != nil {
			logger.Error("Error accepting channel", "err", err)
			continue
		}

		go func(in <-chan *ssh.Request) {
			defer func() {
				_ = ch.Close()
			}()
			for req := range in {
				// TODO check payload
				payload := string(req.Payload)
				switch req.Type {
				case "env":
					var env struct {
						Name  string
						Value string
					}
					if err := ssh.Unmarshal(req.Payload, &env); err != nil {
						logger.Warn("Invalid env payload", "payload", req.Payload, "err", err)
						continue
					}
					if env.Name == "" || env.Value == "" {
						logger.Warn("SSH: Invalid env arguments", "env", env)
						continue
					}

					_, stderr, err := com.ExecCmd("env", fmt.Sprintf("%s=%s", env.Name, env.Value))
					if err != nil {
						logger.Warn("env", "key", err, "value", stderr)
						return
					}

				case "exec":
					cmdName := strings.TrimLeft(payload, "'()")
					logger.Debug("Payload", "name", cmdName)

					args := []string{"serv", "key-" + keyID, "--config=" + conf.CustomConf}
					logger.Debug("Arguments", "args", args)
					cmd := exec.Command(conf.AppPath(), args...)
					cmd.Env = append(os.Environ(), "SSH_ORIGINAL_COMMAND="+cmdName)

					stdout, err := cmd.StdoutPipe()
					if err != nil {
						logger.Error("StdoutPipe", "err", err)
						return
					}
					stderr, err := cmd.StderrPipe()
					if err != nil {
						logger.Error("StderrPipe", "err", err)
						return
					}
					input, err := cmd.StdinPipe()
					if err != nil {
						logger.Error("StdinPipe", "err", err)
						return
					}

					if err = cmd.Start(); err != nil {
						logger.Error("Start", "err", err)
						return
					}

					_ = req.Reply(true, nil)
					go func() {
						_, _ = io.Copy(input, ch)
					}()
					_, _ = io.Copy(ch, stdout)
					_, _ = io.Copy(ch.Stderr(), stderr)

					if err = cmd.Wait(); err != nil {
						logger.Error("Wait", "err", err)
						return
					}

					_, _ = ch.SendRequest("exit-status", false, []byte{0, 0, 0, 0})
					return
				default:
				}
			}
		}(reqs)
	}
}

func setupHostKeys(logger *slog.Logger, appDataPath string, algorithms []string) ([]ssh.Signer, error) {
	dir := filepath.Join(appDataPath, "ssh")
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("create host key directory %w", err)
	}

	var hostKeys []ssh.Signer
	for _, algo := range algorithms {
		keyPath := filepath.Join(dir, "gogs."+algo)
		if !osutil.IsExist(keyPath) {
			args := []string{
				conf.SSH.KeygenPath,
				"-t", algo,
				"-f", keyPath,
				"-m", "PEM",
				"-N", run.Arg(""),
			}
			err = run.Cmd(context.Background(), args...).Run().Wait()
			if err != nil {
				return nil, fmt.Errorf("generate host key with args %v %w", args, err)
			}
			logger.Debug("New private key is generated", "keypath", keyPath)
		}

		keyData, err := os.ReadFile(keyPath)
		if err != nil {
			return nil, fmt.Errorf("read host key %q %w", keyPath, err)
		}
		signer, err := ssh.ParsePrivateKey(keyData)
		if err != nil {
			return nil, fmt.Errorf("parse host key %q %w", keyPath, err)
		}

		hostKeys = append(hostKeys, signer)
	}
	return hostKeys, nil
}
