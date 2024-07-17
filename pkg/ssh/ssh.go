package ssh

import (
	"bytes"
	"io"
	"log/slog"
	"strings"

	"smolgit/pkg/config"
	"smolgit/pkg/git"

	"github.com/go-git/go-billy/v5/osfs"

	"github.com/gliderlabs/ssh"
	"github.com/go-git/go-billy/v5"
	gssh "golang.org/x/crypto/ssh"
)

type Server struct {
	fs  billy.Filesystem
	ssh *ssh.Server
	cfg *config.Config
}

func New(cfg *config.Config) (*Server, error) {
	srv := &Server{
		cfg: cfg,
		fs:  osfs.New(cfg.GitPath),
	}

	srv.ssh = &ssh.Server{
		Handler:          srv.handler,
		PublicKeyHandler: srv.pkHandler,
	}

	return srv, nil
}

func (srv *Server) handler(s ssh.Session) {
	cmd := s.Command()
	slog.Debug("new connection", "cmd", cmd)

	if len(cmd) == 0 {
		cmd = []string{"whoami"}
	}

	var exit int

	switch cmd[0] {
	case "git-receive-pack":
		slog.Debug("receive cmd git-receive-pack", "cmd", cmd)
		exit = srv.cmdRepo(s, cmd)
	case "git-upload-pack":
		slog.Debug("receive cmd git-upload-pack", "cmd", cmd)
		exit = srv.cmdRepo(s, cmd)
	default:
		slog.Debug("command not found\r\n", "cmd", cmd[0])
		exit = 1
	}

	slog.Info("return_code", "exit_code", exit)
	_ = s.Exit(exit)
}

func (srv *Server) pkHandler(ctx ssh.Context, incomingKey ssh.PublicKey) bool {
	slog.Info("handle key", "remote_user", ctx.User(), "remote_addr", ctx.RemoteAddr().String())

	if ctx.User() != "git" {
		slog.Error("wrong remote_user", "user", ctx.User())
		return false
	}

	user, err := srv.cfg.FindUserByKey(string(bytes.TrimSpace(gssh.MarshalAuthorizedKey(incomingKey))))
	if err != nil {
		slog.Error("user not found", "err", err)
		return false
	}
	ctx.SetValue("user_name", user.User)
	slog.Debug("found user", "name", user.User)
	return true
}

func (srv *Server) ListenAndServe() error {
	srv.ssh.Addr = srv.cfg.SSHAddr
	return srv.ssh.ListenAndServe()
}

func (srv *Server) cmdRepo(s ssh.Session, cmd []string) int {
	if len(cmd) != 2 {
		_, _ = io.WriteString(s.Stderr(), "Missing repo name argument\r\n")
		return 1
	}

	repoName := cmd[1]

	userLogin, ok := s.Context().Value("user_name").(string)
	if !ok || userLogin == "" {
		slog.Error("cant find user with", "login", userLogin)
		_, _ = io.WriteString(s.Stderr(), "Permission denied\r\n")
		return 1
	}
	if !strings.HasPrefix(repoName[1:], userLogin) {
		slog.Error("wrong repo prefix", "repoName", repoName, "login", userLogin)
		_, _ = io.WriteString(s.Stderr(), "Permission denied\r\n")
		return 1
	}
	_, err := git.EnsureRepo(srv.fs, srv.cfg.GitBase, repoName)
	if err != nil {
		slog.Error("cant find or create repository", "err", err)
		_, _ = io.WriteString(s.Stderr(), "Repo doesnt exist\r\n")
		return 1
	}
	// TODO sanitize input
	// Get path from user
	slog.Debug("run command", "cmd", cmd, "root", srv.fs.Root(), "path", repoName[1:])
	returnCode := git.RunCommand(srv.fs.Root(), s, []string{cmd[0], repoName[1:]}, []string{})
	return returnCode
}

func Parsepk(data []byte) (ssh.PublicKey, error) {
	publicKey, _, _, _, err := ssh.ParseAuthorizedKey(data) //nolint:dogsled
	return publicKey, err
}
