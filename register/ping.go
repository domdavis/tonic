package register

import (
	"net/http"

	"bitbucket.org/idomdavis/tonic/middleware"
	"github.com/gin-gonic/gin"
)

const ping = "/ping"

// Ping endpoint that can be registered on the router. The endpoint responds
// with pong.
func Ping(r *gin.Engine) {
	r.GET(ping, func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
}

// SilentPing calls Ping, but also adds the endpoint to the skip list on the
// logger.
func SilentPing(r *gin.Engine, reporter *middleware.Reporter) {
	reporter.Skip(ping)
	Ping(r)
}
