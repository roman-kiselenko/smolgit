package route

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"sort"
	"time"

	"smolgit/internal/db"
	"smolgit/internal/git"
	"smolgit/internal/ssh"

	gssh "golang.org/x/crypto/ssh"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	strftime "github.com/itchyny/timefmt-go"

	"github.com/gin-gonic/gin"
)

type Route struct {
	repoLinkBase string
	db           *db.Database
	logger       *slog.Logger
	base         string
	fs           billy.Filesystem
}

//go:embed templates
var htmlTemplates embed.FS

func New(logger *slog.Logger, ginEngine *gin.Engine, database *db.Database, base string) (Route, error) {
	r := Route{
		repoLinkBase: "git@my-git-server.lan",
		db:           database,
		logger:       logger,
		base:         base,
		fs:           osfs.New(base),
	}
	temp := template.New("").Funcs(template.FuncMap{
		"formatAsDate": formatAsDate,
		"formatAsGit":  r.formatAsGit,
	})
	tmpl, err := temp.ParseFS(
		htmlTemplates,
		"templates/layout.html",
		"templates/pages/index.html",
		"templates/pages/500.html",
		"templates/pages/404.html",
		"templates/pages/repo.html",
		"templates/pages/users.html",
		"templates/pages/create.html",
		"templates/pages/files.html",
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
	return fmt.Sprintf("ssh://%s%s", r.repoLinkBase, path)
}

func formatAsDate(t time.Time) string {
	return strftime.Format(t, "%Y/%m/%d %H:%M:%S")
}

func (r *Route) Index(c *gin.Context) {
	repos, _ := r.db.ListRepos()
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Repo",
		"repos": repos,
	})
}

func (r *Route) RepoFiles(c *gin.Context) {
	user, repoPath := c.Param("user"), c.Param("path")
	fullPath := "/" + user + "/" + repoPath
	repo, err := r.db.FindRepoBy(fullPath)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": repoPath + " not found"})
		return
	}
	gitRepo, err := git.OpenRepo(r.logger, r.fs, r.base, fullPath)
	if err != nil {
		r.logger.Error("cant find repository", "err", err)
		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": repoPath + " not found"})
		return
	}
	// TODO consider Until option
	ci, err := gitRepo.Log(&gogit.LogOptions{})
	if err != nil {
		r.logger.Error("go git", "err", err)
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
		r.logger.Error("git repo error", "git", ci)
	}
	defer ci.Close()

	c.HTML(http.StatusOK, "files.html", gin.H{
		"title": "Files",
		"repo":  repo,
		"files": sort.StringSlice(files),
	})
}

func (r *Route) Repo(c *gin.Context) {
	user, repoPath := c.Param("user"), c.Param("path")
	fullPath := "/" + user + "/" + repoPath
	repo, err := r.db.FindRepoBy(fullPath)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": repoPath + " not found"})
		return
	}
	gitRepo, err := git.OpenRepo(r.logger, r.fs, r.base, fullPath)
	if err != nil {
		r.logger.Error("cant find repository", "err", err)
		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": repoPath + " not found"})
		return
	}
	// TODO consider Until option
	ct, err := gitRepo.Log(&gogit.LogOptions{})
	if err != nil {
		r.logger.Error("go git", "err", err)
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
		r.logger.Error("git repo error", "git", ct)
	}
	defer ct.Close()

	c.HTML(http.StatusOK, "repo.html", gin.H{
		"title":   "Repo",
		"repo":    repo,
		"commits": commits,
	})
}

func (r *Route) Users(c *gin.Context) {
	users, _ := r.db.ListUsers()
	c.HTML(http.StatusOK, "users.html", gin.H{
		"title": "Users",
		"users": users,
	})
}

func (r *Route) CreateUser(c *gin.Context) {
	c.HTML(http.StatusOK, "create.html", gin.H{
		"title": "Create Users",
	})
}

func (r *Route) User(c *gin.Context) {
	login := c.PostForm("login")
	password := c.PostForm("password")
	key := c.PostForm("key")
	pk, err := ssh.Parsepk([]byte(key))
	if err != nil {
		r.logger.Error("cant parse key", "err", err)
		c.HTML(http.StatusOK, "500.html", gin.H{
			"title": "Create Users",
			"error": "Cant parse ssh-key",
		})
		return
	}
	kkey := string(bytes.TrimSpace(gssh.MarshalAuthorizedKey(pk)))
	user, _ := r.db.FindUserFromKey(kkey)
	if user.ID != nil && *user.ID > 0 {
		r.logger.Error("cant create user", "err", "duplicate ssh-key")
		c.HTML(http.StatusOK, "500.html", gin.H{
			"title": "Create Users",
			"error": "Duplicate ssh key",
		})
		return
	}
	if err := r.db.InsertUser(login, password, kkey); err != nil {
		r.logger.Error("cant create user", "err", err)
		c.HTML(http.StatusOK, "500.html", gin.H{
			"title": "Create Users",
			"error": "Cant create user",
		})
		return
	}
	c.Redirect(http.StatusFound, "/users")
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
