package route

import (
	"embed"
	_ "embed"
	"html/template"
	"log/slog"
	"net/http"

	"smollgit/internal/db"

	"github.com/gin-gonic/gin"
)

type Route struct {
	db     *db.Database
	logger *slog.Logger
}

//go:embed templates
var htmlTemplates embed.FS

func New(logger *slog.Logger, ginEngine *gin.Engine, database *db.Database) (Route, error) {
	tmpl := template.Must(template.ParseFS(
		htmlTemplates,
		"templates/layout.html",
		"templates/pages/index.html",
		"templates/pages/repos.html",
	))
	ginEngine.SetHTMLTemplate(tmpl)
	return Route{db: database, logger: logger}, nil
}

func (r *Route) Index(c *gin.Context) {
	r.logger.Debug("hit", "route", "index")
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Index",
	})
}

func (r *Route) Repos(c *gin.Context) {
	r.logger.Debug("hit", "route", "repos")
	c.HTML(http.StatusOK, "repos.html", gin.H{
		"title": "Repo",
	})
}

//go:embed templates/css/chota.min.css
var stylesFs embed.FS

func (r *Route) Styles(c *gin.Context) {
	r.logger.Debug("hit", "route", "styles")
	c.Header("Content-Type", "text/css")
	data, _ := stylesFs.ReadFile("templates/css/chota.min.css")
	_, _ = c.Writer.Write(data)
}
