package middleware_test

import (
	"testing"

	"github.com/domdavis/tonic/middleware"
)

func TestLogrus_Log(t *testing.T) {
	t.Run("A nil instance will not panic", func(t *testing.T) {
		t.Parallel()

		l := &middleware.Logrus{}
		l.Log(middleware.ErrorLevel, nil, "test")
	})
}
