package tonic

//nolint:importas // avoiding collision with actual Gin package.
import (
	"fmt"

	"github.com/domdavis/tonic/config"
	"github.com/domdavis/tonic/middleware"
	"github.com/domdavis/tonic/register"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/ulule/limiter/v3"
	limit "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// New returns a preconfigured gin.Engine with logger, limiters, templates
// directory, static content, ping endpoint, among others.
func New(server config.Server, limiters ...config.Limiter) (*gin.Engine, error) {
	if !server.Development {
		gin.SetMode(gin.ReleaseMode)
	}

	logger := middleware.LogrusReporter(logrus.StandardLogger())
	router := gin.New()

	if server.Templates != "" {
		router.LoadHTMLGlob(server.Templates)
	}

	router.Use(GNU)
	router.Use(logger.Log)

	for _, l := range limiters {
		router.Use(limit.NewMiddleware(limiter.New(memory.NewStore(), limiter.Rate{
			Period: l.TTL, Limit: l.Limit})))
	}

	register.SilentPing(router, logger)
	register.Coffee(router)

	if err := register.Static(server.Static, router); err != nil {
		return router, fmt.Errorf("failed to register static content: %w", err)
	}

	return router, nil
}
