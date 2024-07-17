package route

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"sort"
	"strings"
	"time"

	"smolgit/pkg/config"
	"smolgit/pkg/git"
	"smolgit/pkg/model"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	strftime "github.com/itchyny/timefmt-go"

	"github.com/gin-gonic/gin"
)

type Route struct {
	fs  billy.Filesystem
	cfg *config.Config
}

//go:embed templates/html
var htmlTemplates embed.FS

func New(ginEngine *gin.Engine, cfg *config.Config) (Route, error) {
	r := Route{
		fs:  osfs.New(cfg.GitPath),
		cfg: cfg,
	}
	temp := template.New("").Funcs(template.FuncMap{
		"formatAsDate": formatAsDate,
		"formatAsGit":  r.formatAsGit,
		"getBrand":     r.getBrand,
		"formatPath":   formatPath,
	})
	tmpl, err := temp.ParseFS(
		htmlTemplates,
		"templates/html/*.html",
		"templates/html/**/*.html",
		"templates/html/**/**/*.html",
	)
	if err != nil {
		return r, err
	}
	ginEngine.SetHTMLTemplate(tmpl)
	ginEngine.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": "Not Found"})
	})
	return r, nil
}

func (r *Route) formatAsGit(path string) string {
	return fmt.Sprintf("ssh://%s/%s", r.cfg.GitBase, path)
}

func (r *Route) getBrand() string {
	return r.cfg.ServerBrand
}

func formatPath(path string) string {
	chunks := strings.Split(path, "/")
	return strings.TrimSuffix(chunks[2], ".git")
}

func formatAsDate(t time.Time) string {
	return strftime.Format(t, "%Y/%m/%d %H:%M:%S")
}

func (r *Route) Index(c *gin.Context) {
	slog.Debug("hit route", "route", "/")
	repos, err := git.ListRepos(r.cfg.GitPath)
	if err != nil {
		slog.Error("cant find repository", "err", err)
		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": "cant find any repository"})
		return
	}
	for _, repo := range repos {
		gitRepo, _ := git.OpenRepo(r.fs, r.cfg.GitPath, repo.GetFullPath())
		tags, _ := gitRepo.GetTags(true)
		repo.Tags = tags
		refs, _ := gitRepo.GetBranches(true)
		repo.Refs = refs
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Repo",
		"repos": repos,
	})
}

func (r *Route) Refs(c *gin.Context) {
	user, repoPath := c.Param("user"), c.Param("path")
	slog.Debug("hit route", "route", "/repo/refs/:user/:path", "user", user, "path", repoPath)
	fullPath := "/" + user + "/" + repoPath
	// TODO check if repo exist
	gitRepo, err := git.OpenRepo(r.fs, r.cfg.GitPath, fullPath)
	if err != nil {
		slog.Error("cant find repository", "err", err)
		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": repoPath + " not found"})
		return
	}
	tags, err := gitRepo.GetTags(true)
	if err != nil {
		slog.Error("cant get tags", "err", err)
		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": repoPath + " not found"})
		return
	}
	refs, err := gitRepo.GetBranches(true)
	if err != nil {
		slog.Error("cant get branches", "err", err)
		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": repoPath + " not found"})
		return
	}
	c.HTML(http.StatusOK, "refs.html", gin.H{
		"title": "Refs",
		"repo":  model.Repository{Path: repoPath, User: &model.User{Name: user}},
		"refs":  sort.StringSlice(refs),
		"tags":  sort.StringSlice(tags),
	})
}

func (r *Route) Files(c *gin.Context) {
	user, repoPath := c.Param("user"), c.Param("path")
	slog.Debug("hit route", "route", "/repo/files/:user/:path", "user", user, "path", repoPath)
	fullPath := "/" + user + "/" + repoPath
	// TODO check if repo exist
	gitRepo, err := git.OpenRepo(r.fs, r.cfg.GitPath, fullPath)
	if err != nil {
		slog.Error("cant find repository", "err", err)
		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": repoPath + " not found"})
		return
	}
	// TODO consider Until option
	ci, err := gitRepo.Log(&gogit.LogOptions{})
	if err != nil {
		slog.Error("git log", "err", err)
		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": repoPath + " not found"})
		return
	}
	files := []string{}
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
			files = append(files, file.Name)
			return nil
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		slog.Error("repo commits", "err", err)
		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": repoPath + " not found"})
		return
	}
	defer ci.Close()

	c.HTML(http.StatusOK, "files.html", gin.H{
		"title": "Files",
		"repo":  model.Repository{Path: repoPath, User: &model.User{Name: user}},
		"files": sort.StringSlice(files),
	})
}

func (r *Route) Log(c *gin.Context) {
	user, repoPath := c.Param("user"), c.Param("path")
	slog.Debug("hit route", "route", "/repo/log/:user/:path", "user", user, "path", repoPath)
	fullPath := "/" + user + "/" + repoPath
	// TODO check if repo exist
	gitRepo, err := git.OpenRepo(r.fs, r.cfg.GitPath, fullPath)
	if err != nil {
		slog.Error("cant find repository", "err", err)
		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": repoPath + " not found"})
		return
	}
	// TODO consider Until option
	ct, err := gitRepo.Log(&gogit.LogOptions{})
	if err != nil {
		slog.Error("git log", "err", err)
		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": repoPath + " not found"})
		return
	}
	type commit struct {
		Hash    string
		Message string
		Author  string
		Date    string
	}
	commits := []commit{}
	if err := ct.ForEach(func(cmt *object.Commit) error {
		commits = append(commits, commit{
			Hash:    cmt.Hash.String()[0:8],
			Message: cmt.Message,
			Author:  cmt.Author.Name,
			Date:    formatAsDate(cmt.Author.When),
		})
		return nil
	}); err != nil {
		slog.Error("repo commits", "err", err)
		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": repoPath + " not found"})
		return
	}
	defer ct.Close()
	c.HTML(http.StatusOK, "log.html", gin.H{
		"title":   "Repo",
		"repo":    model.Repository{Path: repoPath, User: &model.User{Name: user}},
		"commits": commits,
	})
}

//go:embed templates/css
var styleFs embed.FS

func (r *Route) ExternalStyle(c *gin.Context) {
	c.Header("Content-Type", "text/css")
	data, err := styleFs.ReadFile("templates/css/pico.min.css")
	if err != nil {
		panic(err)
	}
	_, _ = c.Writer.Write(data)
}

func (r *Route) Style(c *gin.Context) {
	c.Header("Content-Type", "text/css")
	data, err := styleFs.ReadFile("templates/css/style.css")
	if err != nil {
		panic(err)
	}
	_, _ = c.Writer.Write(data)
}
