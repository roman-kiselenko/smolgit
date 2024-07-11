package route

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"smolgit/internal/db"

	strftime "github.com/itchyny/timefmt-go"

	"github.com/gin-gonic/gin"
)

type Route struct {
	repoLinkBase string
	db           *db.Database
	logger       *slog.Logger
}

//go:embed templates
var htmlTemplates embed.FS

func New(logger *slog.Logger, ginEngine *gin.Engine, database *db.Database) (Route, error) {
	r := Route{repoLinkBase: "git@my-git-server.lan", db: database, logger: logger}
	temp := template.New("").Funcs(template.FuncMap{
		"formatAsDate": formatAsDate,
		"formatAsGit":  r.formatAsGit,
	})
	tmpl, err := temp.ParseFS(
		htmlTemplates,
		"templates/layout.html",
		"templates/pages/index.html",
		"templates/pages/users.html",
	)
	if err != nil {
		return r, err
	}
	ginEngine.SetHTMLTemplate(tmpl)
	return r, nil
}

func (r *Route) formatAsGit(path string) string {
	return fmt.Sprintf("ssh://%s%s", r.repoLinkBase, path)
}

func formatAsDate(t time.Time) string {
	return strftime.Format(t, "%Y/%m/%d %H:%M:%S")
}

func (r *Route) Index(c *gin.Context) {
	r.logger.Debug("hit", "route", "index")
	repos, err := r.db.ListRepos()
	if err != nil {
		r.logger.Error("cant fetch repos", "err", err)
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Repo",
		"repos": repos,
	})
}

func (r *Route) Users(c *gin.Context) {
	r.logger.Debug("hit", "route", "users")
	users, _ := r.db.ListUsers()
	c.HTML(http.StatusOK, "users.html", gin.H{
		"title": "Users",
		"users": users,
	})
}

//go:embed templates/css
var styleFs embed.FS

func (r *Route) ExternalStyle(c *gin.Context) {
	r.logger.Debug("hit", "route", "style")
	c.Header("Content-Type", "text/css")
	data, err := styleFs.ReadFile("templates/css/terminal.min.css")
	if err != nil {
		r.logger.Warn("cant read css", "err", err)
		c.HTML(http.StatusOK, "repos.html", gin.H{
			"title": "Repo",
		})
		return
	}
	_, _ = c.Writer.Write(data)
}

func (r *Route) Style(c *gin.Context) {
	r.logger.Debug("hit", "route", "style")
	c.Header("Content-Type", "text/css")
	data, err := styleFs.ReadFile("templates/css/style.css")
	if err != nil {
		panic(err)
	}
	_, _ = c.Writer.Write(data)
}
