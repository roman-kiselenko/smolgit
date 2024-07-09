package ssh

import (
	"bytes"
	"context"
	"log/slog"

	"smolgit/internal/db"

	"github.com/gliderlabs/ssh"
	gssh "golang.org/x/crypto/ssh"

	"github.com/urfave/cli/v3"
)

type SSHServer struct {
	addr string

	db     *db.Database
	logger *slog.Logger
	ssh    *ssh.Server
}

func New(logger *slog.Logger, db *db.Database, clictx *cli.Context) (*SSHServer, error) {
	srv := &SSHServer{
		logger: logger,
		db:     db,
		addr:   clictx.String("ssh_addr"),
	}

	if seedRoot(clictx) {
		for _, k := range clictx.StringSlice("root_keys") {
			pk, err := parsepk([]byte(k))
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
		len(clictx.String("root_login")) > 0 &&
		len(clictx.String("root_password")) > 0
}

func (srv *SSHServer) handler(s ssh.Session) {
	var ctx context.Context = s.Context()
	cmd := s.Command()
	srv.logger.Info("new connection", "cmd", cmd)

	if len(cmd) == 0 {
		cmd = []string{"whoami"}
	}

	var exit int

	switch cmd[0] {
	case "git-receive-pack":
		exit = srv.cmdRepo(ctx, s, cmd)
		srv.logger.Debug("git-receive-pack", "code", exit)
		exit = -1
	case "git-upload-pack":
		srv.logger.Debug("git-upload-pack", "code", exit)
		exit = srv.cmdRepo(ctx, s, cmd)
		exit = -1
	default:
		srv.logger.Debug("command not found\r\n", "cmd", cmd[0])
		exit = 1
	}

	srv.logger.Info("return_code", "code", exit)
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
	srv.logger.Debug("found user", "id", *user.ID, "name", user.Login)
	return true
}

func (srv *SSHServer) ListenAndServe() error {
	srv.logger.Info("starting SSH server", "addr", srv.addr)
	srv.ssh.Addr = srv.addr
	return srv.ssh.ListenAndServe()
}

func (srv *SSHServer) cmdRepo(ctx context.Context, s ssh.Session, cmd []string) int {
	if len(cmd) != 2 {
		// _ = writeStringFmt(s.Stderr(), "Missing repo name argument\r\n")
		return 1
	}

	// log, config, user := CtxExtract(ctx)
	// pk := CtxPublicKey(ctx)

	// repoName := sanitize(path.Clean("/" + strings.Trim(cmd[1], "/"))[1:])

	// repo, err := config.LookupRepoAccess(user, repoName)
	// if err != nil {
	// 	_ = writeStringFmt(s.Stderr(), "Repo does not exist\r\n")
	// 	return -1
	// }

	// if repo.Access < access {
	// 	_ = writeStringFmt(s.Stderr(), "Repo does not exist\r\n")
	// 	return -1
	// }

	// if repo.Access >= AccessLevelAdmin {
	// 	_, err = git.EnsureRepo(serv.config.fs, repo.Path())
	// 	if err != nil {
	// 		return -1
	// 	}
	// }

	// returnCode := runCommand(log, serv.fs.Root(), s, []string{cmd[0], repo.Path()}, []string{
	// 	"GITDIR_BASE_DIR=" + serv.fs.Root(),
	// 	"GITDIR_HOOK_REPO_PATH=" + repoName,
	// 	"GITDIR_HOOK_PUBLIC_KEY=" + pk.String(),
	// 	"GITDIR_LOG_FORMAT=console",
	// })

	// if access == AccessLevelWrite {
	// 	switch repo.Type {
	// 	case RepoTypeAdmin, RepoTypeOrgConfig, RepoTypeUserConfig:
	// 		err = serv.Reload()
	// 		if err != nil {
	// 			_ = writeStringFmt(s.Stderr(), "Error when reloading config: %s\r\n", err)
	// 		}
	// 	}
	// }

	return 1
}

func parsepk(data []byte) (ssh.PublicKey, error) {
	publicKey, _, _, _, err := ssh.ParseAuthorizedKey(data)
	return publicKey, err
}
