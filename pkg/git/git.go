package git

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"smolgit/pkg/model"

	"github.com/gliderlabs/ssh"
	billy "github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
)

type Repo struct {
	*git.Repository
}

func RunCommand(
	logger *slog.Logger,
	cwd string,
	session ssh.Session,
	args []string,
	environ []string,
) int {
	cmd := exec.Command(args[0], args[1:]...) //nolint:gosec
	cmd.Dir = cwd

	cmd.Env = append(cmd.Env, environ...)
	cmd.Env = append(cmd.Env, "PATH="+os.Getenv("PATH"))

	stdin, err := cmd.StdinPipe()
	if err != nil {
		logger.Error("failed to get stdin pipe", "err", err)
		return 1
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Error("failed to get stdout pipe", "err", err)
		return 1
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Error("failed to get stderr pipe", "err", err)
		return 1
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	if err = cmd.Start(); err != nil {
		logger.Error("failed to start command", "err", err)
		return 1
	}

	go func() {
		defer stdin.Close()

		if _, stdinErr := io.Copy(stdin, session); stdinErr != nil {
			logger.Error("failed to write session to stdin", "err", err)
		}
	}()

	go func() {
		defer wg.Done()

		if _, stdoutErr := io.Copy(session, stdout); stdoutErr != nil {
			logger.Error("failed to write stdout to session", "err", err)
		}
	}()

	go func() {
		defer wg.Done()

		if _, stderrErr := io.Copy(session.Stderr(), stderr); stderrErr != nil {
			logger.Error("failed to write stderr to session", "err", err)
		}
	}()

	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		logger.Error("failed to wait for command exit", "err", err)
	}

	if err == nil {
		return 0
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		return 1
	}

	return exitErr.ProcessState.ExitCode()
}

func EnsureRepo(logger *slog.Logger, baseFS billy.Filesystem, base, path string) (*Repo, error) {
	info, err := baseFS.Stat(path)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		if err := baseFS.MkdirAll(path, 0o700); err != nil {
			return nil, fmt.Errorf("cant create directory: %s err: %w", path, err)
		}
		logger.Debug("directory created", "path", path)
	}
	fs, err := baseFS.Chroot(path)
	if err != nil {
		return nil, fmt.Errorf("cant chroot to path: %s err: %w", path, err)
	}
	repoFS := filesystem.NewStorage(fs, cache.NewObjectLRUDefault())
	if info == nil {
		logger.Debug("init repo", "path", path)
		repo, err := git.Init(repoFS, nil)
		if err != nil {
			return nil, fmt.Errorf("cant init git: %w", err)
		}
		return &Repo{repo}, nil
	}
	logger.Debug("open repo", "path", path)
	repo, err := git.Open(repoFS, osfs.New(base))
	if err != nil {
		return nil, fmt.Errorf("cant open git: %w", err)
	}

	return &Repo{repo}, nil
}

func ListRepos(base string) ([]model.Repository, error) {
	repos := []model.Repository{}
	info, err := os.Stat(base)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return repos, fmt.Errorf("%s is not exist err: %s", base, err)
	}
	if !info.IsDir() {
		return repos, fmt.Errorf("%s is not a directory", base)
	}
	entries, err := filepath.Glob(base + "/*/*")
	if err != nil {
		return repos, fmt.Errorf("cant read path: %s err: %s", base, err)
	}
	for _, e := range entries {
		repos = append(repos, model.Repository{
			User: &model.User{Login: "admin"},
			Path: strings.TrimPrefix(e, base),
		})
	}
	return repos, nil
}

func OpenRepo(logger *slog.Logger, baseFS billy.Filesystem, base, path string) (*Repo, error) {
	_, err := baseFS.Stat(path)
	if err != nil {
		logger.Debug("cant find path", "path", path, "err", err)
		return nil, fmt.Errorf("cant find path %w %s", err, path)
	}
	fs, err := baseFS.Chroot(path)
	if err != nil {
		return nil, fmt.Errorf("cant chroot to path: %s err: %w", path, err)
	}
	repoFS := filesystem.NewStorage(fs, cache.NewObjectLRUDefault())
	logger.Debug("open repo", "path", path)
	repo, err := git.Open(repoFS, osfs.New(base))
	if err != nil {
		return nil, fmt.Errorf("cant open git: %w", err)
	}

	return &Repo{repo}, nil
}

func (r *Repo) GetTags(clean bool) ([]string, error) {
	ti, err := r.Tags()
	if err != nil {
		return []string{}, err
	}
	tags := []string{}
	if err := ti.ForEach(func(tag *plumbing.Reference) error {
		tags = append(tags, tag.Name().String())
		return nil
	}); err != nil {
		return []string{}, err
	}
	defer ti.Close()
	if clean {
		cleaned := []string{}
		for _, t := range tags {
			cleaned = append(cleaned, strings.TrimPrefix(t, "refs/tags/"))
		}
		return cleaned, nil
	}
	return tags, nil
}

func (r *Repo) GetBranches(clean bool) ([]string, error) {
	bi, err := r.Branches()
	if err != nil {
		return []string{}, err
	}
	refs := []string{}
	if err := bi.ForEach(func(ref *plumbing.Reference) error {
		refs = append(refs, ref.Name().String())
		return nil
	}); err != nil {
		return []string{}, err
	}
	defer bi.Close()
	if clean {
		cleaned := []string{}
		for _, t := range refs {
			cleaned = append(cleaned, strings.TrimPrefix(t, "refs/heads/"))
		}
		return cleaned, nil
	}
	return refs, nil
}
