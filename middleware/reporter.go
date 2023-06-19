package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Reporter is used to report on requests.
type Reporter struct {
	Logger  Logger
	Headers map[string]string
	skip    map[string]struct{}
}

// Logger type used to report the information from the Reporter.
type Logger interface {
	// Log the given message and fields and the correct level.
	Log(level Level, fields map[string]any, message string)
}

// Level used by the reporter to indicate logging level.
type Level uint8

// Log levels.
const (
	ErrorLevel Level = iota
	WarnLevel
	InfoLevel
)

// NewReporter returns a new Reporter for use with Tonic. It will use the given
// Logger to log the reports. A nil logger will cause nothing to be reported.
func NewReporter(logger Logger) *Reporter {
	return &Reporter{
		Logger: logger,
		Headers: map[string]string{
			"X-Forwarded-For":   "forwarded for",
			"X-Forwarded-Proto": "forwarded protocol",
			"X-Forwarded-Port":  "forwarded port",
		},
	}
}

// Skip the given path when logging.
func (r *Reporter) Skip(path string) {
	if r.skip == nil {
		r.skip = map[string]struct{}{}
	}

	r.skip[path] = struct{}{}
}

// Log request the underlying logger.
func (r *Reporter) Log(c *gin.Context) {
	if r.Logger == nil {
		c.Next()

		return
	}

	path := c.Request.URL.Path

	if _, ok := r.skip[path]; ok {
		return
	}

	start := time.Now()

	c.Next()

	latency := Round(time.Since(start))
	statusCode := c.Writer.Status()

	fields := map[string]any{
		"status":  statusCode,
		"method":  c.Request.Method,
		"path":    path,
		"latency": latency,
	}

	fields = r.Embellish(c, fields)

	switch {
	case statusCode >= http.StatusInternalServerError:
		r.Logger.Log(ErrorLevel, fields, "Failed to handle request")
	case statusCode >= http.StatusBadRequest:
		r.Logger.Log(WarnLevel, fields, "Problem handling request")
	case len(c.Errors) > 0:
		r.Logger.Log(ErrorLevel, fields, "Errors handling request")
	default:
		r.Logger.Log(InfoLevel, fields, "Handling request")
	}
}

// Embellish a log message with optional fields if they are set.
func (r *Reporter) Embellish(c *gin.Context, fields map[string]any) map[string]any {
	for header, field := range r.Headers {
		if c.Request.Header.Get(header) != "" {
			fields[field] = c.Request.Header.Get(header)
		}
	}

	if c.ClientIP() != "" {
		fields["client IP"] = c.ClientIP()
	}

	if c.Request.Referer() != "" {
		fields["referer"] = c.Request.Referer()
	}

	if len(c.Errors) > 0 {
		fields["errors"] = c.Errors.JSON()
	}

	return fields
}

// Round a duration to milliseconds if it's over a second.
func Round(duration time.Duration) time.Duration {
	if duration >= time.Second {
		duration = duration.Round(time.Millisecond)
	}

	return duration
}
