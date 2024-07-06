package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Route struct{}

func New() (Route, error) {
	return Route{}, nil
}

func (r *Route) Index(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{})
}
