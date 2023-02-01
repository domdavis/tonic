package register

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Coffee is going to be a problem, we're a teapot.
func Coffee(r *gin.Engine) {
	r.GET("/coffee", func(c *gin.Context) {
		c.String(http.StatusTeapot, http.StatusText(http.StatusTeapot))
	})
}
