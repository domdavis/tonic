package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/domdavis/tonic/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func ExampleReporter_Skip() {
	instance, hook := test.NewNullLogger()
	logger := middleware.LogrusReporter(instance)
	logger.Skip("/skip")

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/logged", nil)

	logger.Log(c)

	fmt.Println(len(hook.Entries))

	c.Request = httptest.NewRequest(http.MethodGet, "/skip", nil)

	logger.Log(c)

	fmt.Println(len(hook.Entries))

	// Output:
	// 1
	// 1
}

func ExampleRound() {
	fmt.Println(middleware.Round(time.Millisecond))
	fmt.Println(middleware.Round(time.Millisecond + time.Nanosecond))
	fmt.Println(middleware.Round(time.Second + time.Millisecond))
	fmt.Println(middleware.Round(time.Second + time.Nanosecond))

	// Output:
	// 1ms
	// 1.000001ms
	// 1.001s
	// 1s
}

func TestLogger_Log(t *testing.T) {
	t.Run("Normal log messages are info level", func(t *testing.T) {
		t.Parallel()

		instance, hook := test.NewNullLogger()
		logger := middleware.LogrusReporter(instance)
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/logged", nil)

		logger.Log(c)

		assert.Equal(t, hook.LastEntry().Level, logrus.InfoLevel)
	})

	t.Run("Context errors are logged", func(t *testing.T) {
		t.Parallel()

		instance, hook := test.NewNullLogger()
		logger := middleware.LogrusReporter(instance)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/logged", nil)
		c.Errors = []*gin.Error{{Err: assert.AnError}}
		logger.Log(c)

		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	})

	t.Run("500 codes are error level", func(t *testing.T) {
		t.Parallel()

		instance, hook := test.NewNullLogger()
		logger := middleware.LogrusReporter(instance)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Request = httptest.NewRequest(http.MethodGet, "/logged", nil)

		logger.Log(c)

		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	})

	t.Run("400 codes are a warning", func(t *testing.T) {
		t.Parallel()

		instance, hook := test.NewNullLogger()
		logger := middleware.LogrusReporter(instance)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.Request = httptest.NewRequest(http.MethodGet, "/logged", nil)

		logger.Log(c)

		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("Setting no logger will not panic", func(t *testing.T) {
		t.Parallel()

		logger := middleware.NewReporter(nil)
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/logged", nil)

		logger.Log(c)
	})
}

func TestLogger_Embellish(t *testing.T) {
	t.Run("Optional elements are excluded if absent", func(t *testing.T) {
		t.Parallel()

		instance, hook := test.NewNullLogger()
		logger := middleware.LogrusReporter(instance)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/logged", nil)
		c.Request.RemoteAddr = ""

		logger.Log(c)

		for _, entry := range logger.Headers {
			assert.NotContains(t, hook.LastEntry().Data, entry)
		}

		assert.NotContains(t, hook.LastEntry().Data, "client IP")
		assert.NotContains(t, hook.LastEntry().Data, "referer")
		assert.NotContains(t, hook.LastEntry().Data, "errors")
	})

	t.Run("Optional elements are added if present", func(t *testing.T) {
		t.Parallel()

		instance, hook := test.NewNullLogger()
		logger := middleware.LogrusReporter(instance)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/logged", nil)
		c.Errors = []*gin.Error{{Err: assert.AnError}}
		c.Request.RemoteAddr = "127.0.0.1:8000"
		c.Request.Header["Referer"] = []string{"test"}

		for k, v := range logger.Headers {
			c.Request.Header[k] = []string{v}
		}

		logger.Log(c)

		for _, entry := range logger.Headers {
			assert.Contains(t, hook.LastEntry().Data, entry)
		}

		assert.Contains(t, hook.LastEntry().Data, "client IP")
		assert.Contains(t, hook.LastEntry().Data, "referer")
		assert.Contains(t, hook.LastEntry().Data, "errors")
	})
}
