package tonic

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logger type for use with Gin.
type Logger struct {
	Instance *logrus.Logger
	Headers  map[string]string
	skip     map[string]struct{}
}

// NewLogger returns a new Logger instance using the Logrus standard logger and
// some default headers to log.
func NewLogger(logger *logrus.Logger) *Logger {
	return &Logger{
		Instance: logger,
		Headers: map[string]string{
			"X-Forwarded-For":   "forwarded for",
			"X-Forwarded-Proto": "forwarded protocol",
			"X-Forwarded-Port":  "forwarded port",
		},
	}
}

// Skip the given path when logging.
func (l *Logger) Skip(path string) {
	if l.skip == nil {
		l.skip = map[string]struct{}{}
	}

	l.skip[path] = struct{}{}
}

// Log request using Logrus. If no logrus Logger is set then the standard
// Logger is used.
func (l *Logger) Log(c *gin.Context) {
	path := c.Request.URL.Path

	if _, ok := l.skip[path]; ok {
		return
	}

	start := time.Now()

	c.Next()

	latency := Round(time.Since(start))
	statusCode := c.Writer.Status()

	fields := logrus.Fields{
		"status":  statusCode,
		"method":  c.Request.Method,
		"path":    path,
		"latency": latency,
	}

	l.Embellish(c, fields)

	if l.Instance == nil {
		l.Instance = logrus.StandardLogger()
	}

	switch {
	case statusCode >= http.StatusInternalServerError:
		l.Instance.WithFields(fields).Error("Failed to handle request")
	case statusCode >= http.StatusBadRequest:
		l.Instance.WithFields(fields).Warn("Problem handling request")
	case len(c.Errors) > 0:
		l.Instance.WithFields(fields).Error("Errors handling request")
	default:
		l.Instance.WithFields(fields).Info("Handling request")
	}
}

// Embellish a log message with optional fields if they are set.
func (l *Logger) Embellish(c *gin.Context, fields logrus.Fields) logrus.Fields {
	for header, field := range l.Headers {
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
