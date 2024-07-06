package route

import (
	"net/http"

	"git-pebble/internal/db"

	"github.com/gin-gonic/gin"
)

type Route struct {
	db *db.Database
}

func New(db *db.Database) (Route, error) {
	return Route{
		db: db,
	}, nil
}

func (r *Route) Index(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{})
}
