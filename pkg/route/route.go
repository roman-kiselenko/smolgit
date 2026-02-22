package route

import (
	"log/slog"
	"net/http"

	"smolgit/pkg/config"
	"smolgit/pkg/git"
	"smolgit/pkg/model"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"

	"github.com/gin-gonic/gin"
)

type Route struct {
	fs    billy.Filesystem
	cfg   *config.Config
	users []model.User
}

func New(ginEngine *gin.Engine, cfg *config.Config) (Route, error) {
	r := Route{
		fs:    osfs.New(cfg.GitPath),
		cfg:   cfg,
		users: cfg.Users,
	}

	return r, nil
}

func (r *Route) Repos(c *gin.Context) {
	slog.Debug("hit route", "route", "/")
	repos, err := git.ListRepos(r.cfg.GitPath)
	if err != nil {
		slog.Error("cant find repository", "err", err)
		c.JSON(http.StatusNotFound, gin.H{"title": "cant find any repository"})
		return
	}
	for _, repo := range repos {
		gitRepo, _ := git.OpenRepo(r.fs, r.cfg.GitPath, repo.GetFullPath())
		tags, _ := gitRepo.GetTags(true)
		repo.Tags = tags
		refs, _ := gitRepo.GetBranches(true)
		repo.Refs = refs
	}

	c.JSON(http.StatusOK, gin.H{
		"meta": map[string]interface{}{
			"total": len(repos),
		},
		"items": repos,
	})
}

func (r *Route) Commit(c *gin.Context) {
	user, repoPath := c.Param("user"), c.Param("path")
	slog.Debug("hit route", "route", "/repo/commit/:user/:path", "user", user, "path", repoPath)
	c.JSON(http.StatusOK, gin.H{"title": "Commit", "repo": model.Repository{Path: repoPath, User: &model.User{Name: user}}})
}
