package route

import (
	"log/slog"
	"net/http"

	"smollgit/internal/db"

	"github.com/gin-gonic/gin"
)

type Route struct {
	db     *db.Database
	logger *slog.Logger
}

func New(logger *slog.Logger, database *db.Database) (Route, error) {
	return Route{db: database, logger: logger}, nil
}

func (r *Route) Index(c *gin.Context) {
	r.logger.Debug("hit", "route", "index")
	c.JSON(http.StatusOK, map[string]interface{}{})
}
