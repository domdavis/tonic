package jwt_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bitbucket.org/idomdavis/tonic"
	"bitbucket.org/idomdavis/tonic/config"
	"bitbucket.org/idomdavis/tonic/jwt"
	"bitbucket.org/idomdavis/tonic/middleware"
	"bitbucket.org/idomdavis/tonic/register"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func ExampleSign() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := config.Security{Secret: "secret", SessionTTL: time.Hour}

	// Set a JWT with claims from the "user" and "realm" params of the
	// context.
	jwt.Sign(c, s, "user", "realm")

	b, _ := io.ReadAll(w.Result().Body)

	fmt.Println(w.Result().StatusCode, len(b) > 0)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	// Force the secret to be random.
	s.Secret = ""

	// Set a JWT with no claims.
	jwt.Sign(c, s)

	b, _ = io.ReadAll(w.Result().Body)

	fmt.Println(w.Result().StatusCode, len(b) > 0)

	// Output:
	// 200 true
	// 200 true
}

func ExampleSet() {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	jwt.Set(req, "token")
	fmt.Println(req.Header.Get(jwt.Header))

	// Output:
	// Bearer token
}

func TestSign(t *testing.T) {
	t.Run("A zero TTL will error", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		s := config.Security{}

		jwt.Sign(c, s)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		assert.NotEmpty(t, c.Errors)
	})

	t.Run("A negative TTL will error", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		s := config.Security{SessionTTL: -1}

		jwt.Sign(c, s)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		assert.NotEmpty(t, c.Errors)
	})

	t.Run("Invalid claims will error", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("chan", make(chan int))
		s := config.Security{SessionTTL: time.Hour}

		jwt.Sign(c, s, "chan")

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		assert.NotEmpty(t, c.Errors)
	})
}

func TestAuthenticate(t *testing.T) {
	const endpoint = "/ping"

	t.Run("A valid token will authenticate", func(t *testing.T) {
		t.Parallel()

		security := config.Security{Secret: "secret", SessionTTL: time.Hour}

		logger := middleware.LogrusReporter(logrus.StandardLogger())
		logger.Skip(endpoint)

		router := gin.New()
		router.Use(logger.Log)
		router.Use(jwt.Authenticate(security))
		register.Ping(router)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)

		assert.NoError(t, err)

		signatory := &tonic.Signatory{TTL: security.SessionTTL, Secret: security.Secret}
		signatory.Initialise()

		token, err := signatory.Sign(&gin.Context{})

		assert.NoError(t, err)

		jwt.Set(req, token)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("An invalid token will fail", func(t *testing.T) {
		t.Parallel()

		security := config.Security{Secret: "secret", SessionTTL: time.Hour}

		logger := middleware.LogrusReporter(logrus.StandardLogger())
		logger.Skip(endpoint)

		router := gin.New()
		router.Use(logger.Log)
		router.Use(jwt.Authenticate(security))
		register.Ping(router)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)

		assert.NoError(t, err)

		jwt.Set(req, "token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})
}
