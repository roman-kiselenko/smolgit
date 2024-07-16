package ssh

import (
	"bytes"
	"io"
	"log/slog"
	"strings"

	"smolgit/internal/db"
	"smolgit/internal/git"

	"github.com/go-git/go-billy/v5/osfs"

	"github.com/gliderlabs/ssh"
	"github.com/go-git/go-billy/v5"
	gssh "golang.org/x/crypto/ssh"

	"github.com/urfave/cli/v3"
)

type SSHServer struct {
	addr string

	base   string
	fs     billy.Filesystem
	db     *db.Database
	logger *slog.Logger
	ssh    *ssh.Server
}

func New(logger *slog.Logger, db *db.Database, clictx *cli.Context) (*SSHServer, error) {
	srv := &SSHServer{
		logger: logger,
		db:     db,
		fs:     osfs.New(clictx.String("git_path")),
		base:   clictx.String("git_path"),
		addr:   clictx.String("ssh_addr"),
	}

	if seedRoot(clictx) {
		for _, k := range clictx.StringSlice("root_keys") {
			pk, err := Parsepk([]byte(k))
			if err != nil {
				srv.logger.Error("cant parse ssh_keys", "key", k)
				continue
			}
			if root, err := db.FindRootUser(); err == nil && root.ID != nil {
				if err := db.InsertUserKey(*root.ID, string(bytes.TrimSpace(gssh.MarshalAuthorizedKey(pk)))); err == nil {
					srv.logger.Info("ssh_keys created", "key", "*****")
				}
			}
		}
	}

	srv.ssh = &ssh.Server{
		Handler:          srv.handler,
		PublicKeyHandler: srv.pkHandler,
	}

	return srv, nil
}

func seedRoot(clictx *cli.Context) bool {
	return len(clictx.StringSlice("root_keys")) > 0 &&
		clictx.String("root_login") != "" &&
		clictx.String("root_password") != ""
}

func (srv *SSHServer) handler(s ssh.Session) {
	cmd := s.Command()
	srv.logger.Info("new connection", "cmd", cmd)

	if len(cmd) == 0 {
		cmd = []string{"whoami"}
	}

	var exit int

	switch cmd[0] {
	case "git-receive-pack":
		srv.logger.Debug("receive cmd git-receive-pack", "cmd", cmd)
		exit = srv.cmdRepo(s, cmd)
	case "git-upload-pack":
		srv.logger.Debug("receive cmd git-upload-pack", "cmd", cmd)
		exit = srv.cmdRepo(s, cmd)
	default:
		srv.logger.Debug("command not found\r\n", "cmd", cmd[0])
		exit = 1
	}

	srv.logger.Info("return_code", "exit_code", exit)
	_ = s.Exit(exit)
}

func (srv *SSHServer) pkHandler(ctx ssh.Context, incomingKey ssh.PublicKey) bool {
	srv.logger.Info("handle key", "remote_user", ctx.User(), "remote_addr", ctx.RemoteAddr().String())

	if ctx.User() != "git" {
		srv.logger.Error("wrong remote_user", "user", ctx.User())
		return false
	}

	user, err := srv.db.FindUserFromKey(string(bytes.TrimSpace(gssh.MarshalAuthorizedKey(incomingKey))))
	if err != nil {
		srv.logger.Error("user not found", "err", err)
		return false
	}
	ctx.SetValue("user_id", *user.ID)
	ctx.SetValue("user_name", user.Login)
	srv.logger.Debug("found user", "id", *user.ID, "name", user.Login)
	return true
}

func (srv *SSHServer) ListenAndServe() error {
	srv.ssh.Addr = srv.addr
	return srv.ssh.ListenAndServe()
}

func (srv *SSHServer) cmdRepo(s ssh.Session, cmd []string) int {
	if len(cmd) != 2 {
		_, _ = io.WriteString(s.Stderr(), "Missing repo name argument\r\n")
		return 1
	}

	repoName := cmd[1]

	// TODO dry
	userID, ok := s.Context().Value("user_id").(int64)
	if !ok || userID == 0 {
		srv.logger.Error("cant find user with", "id", userID)
		_, _ = io.WriteString(s.Stderr(), "Permission denied\r\n")
		return 1
	}
	userLogin, ok := s.Context().Value("user_name").(string)
	if !ok || userLogin == "" {
		srv.logger.Error("cant find user with", "login", userLogin)
		_, _ = io.WriteString(s.Stderr(), "Permission denied\r\n")
		return 1
	}
	if !strings.HasPrefix(repoName[1:], userLogin) {
		srv.logger.Error("wrong repo prefix", "repoName", repoName, "login", userLogin)
		_, _ = io.WriteString(s.Stderr(), "Permission denied\r\n")
		return 1
	}
	_, err := git.EnsureRepo(srv.logger, srv.fs, srv.base, repoName)
	if err != nil {
		srv.logger.Error("cant find or create repository", "err", err)
		_, _ = io.WriteString(s.Stderr(), "Repo doesnt exist\r\n")
		return 1
	}
	if !srv.db.RepoExist(userID, repoName) {
		if err := srv.db.InsertRepo(userID, repoName); err != nil {
			srv.logger.Error("cant insert repo", "err", err)
			return 1
		}
	}
	// TODO sanitize input
	// Get path from user
	srv.logger.Debug("run command", "cmd", cmd, "root", srv.fs.Root(), "path", repoName[1:])
	returnCode := git.RunCommand(srv.logger, srv.fs.Root(), s, []string{cmd[0], repoName[1:]}, []string{})
	return returnCode
}

func Parsepk(data []byte) (ssh.PublicKey, error) {
	publicKey, _, _, _, err := ssh.ParseAuthorizedKey(data) //nolint:dogsled
	return publicKey, err
}
