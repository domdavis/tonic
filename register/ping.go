package register

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping endpoint that can be registered on the router. The endpoint responds
// with pong.
func Ping(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
}
