package ssh

import (
	"bytes"
	"log/slog"

	"smolgit/internal/db"

	"github.com/gliderlabs/ssh"
	gssh "golang.org/x/crypto/ssh"

	"github.com/urfave/cli/v3"
)

type SSHServer struct {
	addr string

	logger *slog.Logger
	ssh    *ssh.Server
}

func New(logger *slog.Logger, db *db.Database, clictx *cli.Context) (*SSHServer, error) {
	srv := &SSHServer{
		logger: logger,
		addr:   clictx.String("ssh_addr"),
	}

	srv.ssh = &ssh.Server{
		Handler: func(s ssh.Session) {
			cmd := s.Command()
			logger.Info("new connection", "cmd", cmd)

			if len(cmd) == 0 {
				cmd = []string{"whoami"}
			}

			var exit int

			switch cmd[0] {
			case "whoami":
				// _ = writeStringFmt(s, "logged in as %s\r\n", user.Username)
				// exit = cmdWhoami(ctx, s, cmd)
				exit = -1
			case "git-receive-pack":
				// exit = srv.cmdRepoAction(ctx, s, cmd, AccessLevelWrite)
				exit = -1
			case "git-upload-pack":
				// exit = srv.cmdRepoAction(ctx, s, cmd, AccessLevelRead)
				exit = -1
			default:
				logger.Debug("command not found\r\n", "cmd", cmd[0])
				exit = 1
			}

			logger.Info("return_code", "code", exit)
			_ = s.Exit(exit)
		},
		PublicKeyHandler: func(ctx ssh.Context, incomingKey ssh.PublicKey) bool {
			logger.Info("handle key", "remote_user", ctx.User(), "remote_addr", ctx.RemoteAddr().String())

			remoteUser := ctx.User()

			user, err := db.FindUserFromKey(
				string(bytes.TrimSpace(gssh.MarshalAuthorizedKey(incomingKey))),
				remoteUser,
			)
			if err != nil {
				logger.Error("user not found", "err", err)
				return false
			}
			logger.Debug("found user", "user", user)
			return true
		},
	}

	return srv, nil
}

func (srv *SSHServer) ListenAndServe() error {
	srv.logger.Info("starting SSH server", "addr", srv.addr)
	srv.ssh.Addr = srv.addr
	return srv.ssh.ListenAndServe()
}
