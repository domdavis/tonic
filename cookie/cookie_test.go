package cookie_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/domdavis/tonic"
	"github.com/domdavis/tonic/config"
	"github.com/domdavis/tonic/cookie"
	"github.com/domdavis/tonic/middleware"
	"github.com/domdavis/tonic/register"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func ExampleDrop() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := config.Security{Secret: "secret", SessionTTL: time.Hour}

	// Drop a cookie with claims from the "user" and "realm" params of the
	// context.
	fmt.Println(cookie.Drop(c, s, "user", "realm"))
	fmt.Println(len(w.Result().Cookies()))

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	// Force the secret to be random.
	s.Secret = ""

	// Drop a cookie with no claims.
	// context.
	fmt.Println(cookie.Drop(c, s))
	fmt.Println(len(w.Result().Cookies()))

	// Output:
	// <nil>
	// 1
	// <nil>
	// 1
}

func ExampleClear() {
	s := config.Security{Secret: "secret", SessionTTL: time.Hour}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	cookie.Clear(c, s)

	fmt.Println(w.Result().Cookies()[0].MaxAge)

	// Output:
	// -1
}

func TestDrop(t *testing.T) {
	t.Run("A zero TTL will error", func(t *testing.T) {
		t.Parallel()

		c := &gin.Context{}
		s := config.Security{}

		err := cookie.Drop(c, s)

		assert.ErrorIs(t, err, tonic.ErrInvalidTTL)
	})

	t.Run("A negative TTL will error", func(t *testing.T) {
		t.Parallel()

		c := &gin.Context{}
		s := config.Security{SessionTTL: -1}

		err := cookie.Drop(c, s)

		assert.ErrorIs(t, err, tonic.ErrInvalidTTL)
	})

	t.Run("A sub-second will error", func(t *testing.T) {
		t.Parallel()

		c := &gin.Context{}
		s := config.Security{SessionTTL: 1}

		err := cookie.Drop(c, s)

		assert.ErrorIs(t, err, tonic.ErrInvalidTTL)
	})

	t.Run("Invalid claims will error", func(t *testing.T) {
		t.Parallel()

		c := &gin.Context{}
		c.Set("chan", make(chan int))
		s := config.Security{SessionTTL: time.Hour}

		err := cookie.Drop(c, s, "chan")

		assert.ErrorContains(t, err, "failed to drop")
	})
}

func TestAuthenticate(t *testing.T) {
	const endpoint = "/ping"

	security := config.Security{Secret: "secret", SessionTTL: time.Hour}

	t.Run("A valid cookie will authenticate", func(t *testing.T) {
		t.Parallel()

		logger := middleware.LogrusReporter(logrus.StandardLogger())
		logger.Skip(endpoint)

		router := gin.New()
		router.Use(logger.Log)
		router.Use(cookie.Authenticate(security, ""))

		register.Ping(router)

		signatory := &tonic.Signatory{TTL: security.SessionTTL, Secret: security.Secret}
		signatory.Initialise()

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)

		assert.NoError(t, err)

		token, err := signatory.Sign(&gin.Context{})

		assert.NoError(t, err)

		req.AddCookie(&http.Cookie{
			Name:     cookie.Name,
			Value:    token,
			Path:     "/",
			MaxAge:   100,
			Secure:   true,
			HttpOnly: true,
		})

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("An invalid cookie will fail", func(t *testing.T) {
		t.Parallel()

		logger := middleware.LogrusReporter(logrus.StandardLogger())
		logger.Skip(endpoint)

		router := gin.New()
		router.Use(logger.Log)
		router.Use(cookie.Authenticate(security, ""))

		register.Ping(router)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)

		assert.NoError(t, err)

		req.AddCookie(&http.Cookie{
			Name:     cookie.Name,
			Value:    "garbage",
			Path:     "/",
			MaxAge:   100,
			Secure:   true,
			HttpOnly: true,
		})

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})

	t.Run("A missing cookie will fail", func(t *testing.T) {
		t.Parallel()

		logger := middleware.LogrusReporter(logrus.StandardLogger())
		logger.Skip(endpoint)

		router := gin.New()
		router.Use(logger.Log)
		router.Use(cookie.Authenticate(security, ""))

		register.Ping(router)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)

		assert.NoError(t, err)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})

	t.Run("Redirects will occur on auth failure", func(t *testing.T) {
		t.Parallel()

		logger := middleware.LogrusReporter(logrus.StandardLogger())
		logger.Skip(endpoint)

		router := gin.New()
		router.Use(logger.Log)
		router.Use(cookie.Authenticate(security, "/redirect"))

		register.Ping(router)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)

		assert.NoError(t, err)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Result().StatusCode)
	})
}
