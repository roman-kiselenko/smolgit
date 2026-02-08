package route

import (
	"log/slog"
	"net/http"
	"slices"
	"time"

	"smolgit/pkg/config"
	"smolgit/pkg/git"
	"smolgit/pkg/model"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

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

type creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *Route) Login(c *gin.Context) {
	var req creds
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("parsing", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	index := slices.IndexFunc(r.users, func(u model.User) bool {
		return u.Name == req.Username
	})
	if index == -1 {
		slog.Error("no such user", "user", req.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}
	u := r.users[index]
	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)) != nil {
		slog.Error("bad password", "user", req.Username, "password", req.Password)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	exp := time.Now().Add(1 * time.Hour)
	claims := &model.Claims{
		Username: u.Name,
		Role:     u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(r.cfg.ServerJWTKey))
	if err != nil {
		slog.Error("cant sign token", "user", req.Username, "password", req.Password)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": t})
}

func (r *Route) Files(c *gin.Context) {
	user, repoPath := c.Param("user"), c.Param("path")
	slog.Debug("hit route", "route", "/repo/files/:user/:path", "user", user, "path", repoPath)
	fullPath := "/" + user + "/" + repoPath
	gitRepo, err := git.OpenRepo(r.fs, r.cfg.GitPath, fullPath)
	if err != nil {
		slog.Error("cant find repository", "err", err)
		c.JSON(http.StatusNotFound, gin.H{"title": "cant find any repository"})
		return
	}
	// TODO consider Until option
	ci, err := gitRepo.Log(&gogit.LogOptions{})
	if err != nil {
		slog.Error("git log", "err", err)
		c.JSON(http.StatusNotFound, gin.H{"title": repoPath + " not found"})
		return
	}
	files := []map[string]string{}
	count := 0
	if err := ci.ForEach(func(cmt *object.Commit) error {
		if count == 1 {
			return nil
		}
		fi, err := cmt.Files()
		if err != nil {
			return err
		}
		count++
		if err := fi.ForEach(func(file *object.File) error {
			files = append(files, map[string]string{"filename": file.Name})
			return nil
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		slog.Error("repo commits", "err", err)
		c.JSON(http.StatusNotFound, gin.H{"title": repoPath + " not found"})
		return
	}
	defer ci.Close()

	c.JSON(http.StatusOK, gin.H{
		"title": "Files",
		"repo":  model.Repository{Path: repoPath, User: &model.User{Name: user}},
		"files": files,
	})
}

func (r *Route) Commit(c *gin.Context) {
	user, repoPath := c.Param("user"), c.Param("path")
	slog.Debug("hit route", "route", "/repo/commit/:user/:path", "user", user, "path", repoPath)
	c.JSON(http.StatusOK, gin.H{"title": "Commit", "repo": model.Repository{Path: repoPath, User: &model.User{Name: user}}})
}
