package route

import (
	"log/slog"
	"net/http"
	"time"

	"smolgit/pkg/config"
	"smolgit/pkg/git"
	"smolgit/pkg/model"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

type Route struct {
	fs    billy.Filesystem
	cfg   *config.Config
	users *config.Users
}

func New(ginEngine *gin.Engine, cfg *config.Config) (Route, error) {
	r := Route{
		fs:  osfs.New(cfg.GitPath),
		cfg: cfg,
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
		"title": "Repo",
		"repos": repos,
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
	u, ok := r.users.Users[req.Username]
	if !ok || bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	exp := time.Now().Add(1 * time.Hour)
	claims := &model.Claims{
		Username: u.Username,
		Role:     u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(r.cfg.ServerJWTKey))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": t})
}

// func (r *Route) Refs(c *gin.Context) {
// 	user, repoPath := c.Param("user"), c.Param("path")
// 	slog.Debug("hit route", "route", "/repo/refs/:user/:path", "user", user, "path", repoPath)
// 	fullPath := "/" + user + "/" + repoPath
// 	// TODO check if repo exist
// 	gitRepo, err := git.OpenRepo(r.fs, r.cfg.GitPath, fullPath)
// 	if err != nil {
// 		slog.Error("cant find repository", "err", err)
// 		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": repoPath + " not found"})
// 		return
// 	}
// 	tags, err := gitRepo.GetTags(true)
// 	if err != nil {
// 		slog.Error("cant get tags", "err", err)
// 		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": repoPath + " not found"})
// 		return
// 	}
// 	refs, err := gitRepo.GetBranches(true)
// 	if err != nil {
// 		slog.Error("cant get branches", "err", err)
// 		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": repoPath + " not found"})
// 		return
// 	}
// 	c.HTML(http.StatusOK, "refs.html", gin.H{
// 		"title": "Refs",
// 		"repo":  model.Repository{Path: repoPath, User: &model.User{Name: user}},
// 		"refs":  sort.StringSlice(refs),
// 		"tags":  sort.StringSlice(tags),
// 	})
// }

// func (r *Route) Files(c *gin.Context) {
// 	user, repoPath := c.Param("user"), c.Param("path")
// 	slog.Debug("hit route", "route", "/repo/files/:user/:path", "user", user, "path", repoPath)
// 	fullPath := "/" + user + "/" + repoPath
// 	// TODO check if repo exist
// 	gitRepo, err := git.OpenRepo(r.fs, r.cfg.GitPath, fullPath)
// 	if err != nil {
// 		slog.Error("cant find repository", "err", err)
// 		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": repoPath + " not found"})
// 		return
// 	}
// 	// TODO consider Until option
// 	ci, err := gitRepo.Log(&gogit.LogOptions{})
// 	if err != nil {
// 		slog.Error("git log", "err", err)
// 		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": repoPath + " not found"})
// 		return
// 	}
// 	files := []string{}
// 	count := 0
// 	if err := ci.ForEach(func(cmt *object.Commit) error {
// 		if count == 1 {
// 			return nil
// 		}
// 		fi, err := cmt.Files()
// 		if err != nil {
// 			return err
// 		}
// 		count++
// 		if err := fi.ForEach(func(file *object.File) error {
// 			files = append(files, file.Name)
// 			return nil
// 		}); err != nil {
// 			return err
// 		}
// 		return nil
// 	}); err != nil {
// 		slog.Error("repo commits", "err", err)
// 		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": repoPath + " not found"})
// 		return
// 	}
// 	defer ci.Close()

// 	c.HTML(http.StatusOK, "files.html", gin.H{
// 		"title": "Files",
// 		"repo":  model.Repository{Path: repoPath, User: &model.User{Name: user}},
// 		"files": sort.StringSlice(files),
// 	})
// }

// func (r *Route) Log(c *gin.Context) {
// 	user, repoPath := c.Param("user"), c.Param("path")
// 	slog.Debug("hit route", "route", "/repo/log/:user/:path", "user", user, "path", repoPath)
// 	fullPath := "/" + user + "/" + repoPath
// 	// TODO check if repo exist
// 	gitRepo, err := git.OpenRepo(r.fs, r.cfg.GitPath, fullPath)
// 	if err != nil {
// 		slog.Error("cant find repository", "err", err)
// 		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": repoPath + " not found"})
// 		return
// 	}
// 	// TODO consider Until option
// 	ct, err := gitRepo.Log(&gogit.LogOptions{})
// 	if err != nil {
// 		slog.Error("git log", "err", err)
// 		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": repoPath + " not found"})
// 		return
// 	}
// 	type commit struct {
// 		Hash    string
// 		Message string
// 		Author  string
// 		Date    string
// 	}
// 	commits := []commit{}
// 	if err := ct.ForEach(func(cmt *object.Commit) error {
// 		commits = append(commits, commit{
// 			Hash:    cmt.Hash.String()[0:8],
// 			Message: cmt.Message,
// 			Author:  cmt.Author.Name,
// 		})
// 		return nil
// 	}); err != nil {
// 		slog.Error("repo commits", "err", err)
// 		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": repoPath + " not found"})
// 		return
// 	}
// 	defer ct.Close()
// 	c.HTML(http.StatusOK, "log.html", gin.H{
// 		"title":   "Repo",
// 		"repo":    model.Repository{Path: repoPath, User: &model.User{Name: user}},
// 		"commits": commits,
// 	})
// }
