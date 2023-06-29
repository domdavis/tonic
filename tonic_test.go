package tonic_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	"github.com/domdavis/tonic"
	"github.com/domdavis/tonic/config"
	"github.com/domdavis/tonic/cookie"
	"github.com/domdavis/tonic/jwt"
	"github.com/gin-gonic/gin"
)

var Config = struct {
	config.Server
	config.Security
	config.Limiter
}{
	Server: config.Server{
		Templates:   "testdata/templates/*",
		Static:      "testdata/static",
		Development: true,
	},
	Security: config.Security{
		Secret:     "secret",
		SessionTTL: time.Second * 60,
		Timebox:    100,
	},
	Limiter: config.Limiter{
		Limit: 100,
		TTL:   60,
	},
}

//nolint:misspell // HTTP status codes are in American English.
func Example_tonic() {
	router, err := tonic.New(Config.Server, Config.Limiter)

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	router.GET("/login", Login)
	router.GET("/token", Token)

	webapp := router.Group("/webapp")
	webapp.Use(cookie.Authenticate(Config.Security, "/webapp/login"))
	webapp.GET("/page", Endpoint)

	api := router.Group("/api")
	api.Use(jwt.Authenticate(Config.Security))
	api.GET("/action", Endpoint)

	r := call(router, "/ping", nil)
	r = call(router, "/webapp/page", r.Result())
	r = call(router, "/api/action", r.Result())

	r = call(router, "/login", nil)

	call(router, "/webapp/page", r.Result())

	r = call(router, "/token", nil)

	call(router, "/api/action", r.Result())

	// Output:
	// 200 pong
	// 307 <a href="/webapp/login">Temporary Redirect</a>.
	// 401 Unauthorized
	// 200 OK
	// 200 OK
	// 200 Bearer
	// 200 OK
}

func call(router *gin.Engine, endpoint string, response *http.Response) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)

	if err != nil {
		fmt.Println(err)

		return w
	}

	if response != nil {
		for _, c := range response.Cookies() {
			req.AddCookie(c)
		}

		if t := response.Header.Get("Content-Type"); t == jwt.ContentType {
			b, _ := io.ReadAll(response.Body)

			jwt.Set(req, string(b))
		}
	}

	router.ServeHTTP(w, req)

	if t := w.Header().Get("Content-Type"); t == jwt.ContentType {
		fmt.Println(w.Code, jwt.Prefix)
	} else {
		fmt.Println(w.Code, strings.TrimSpace(w.Body.String()))
	}

	return w
}

// Login using a constant time algorithm. No actual authentication is done for
// this example.
func Login(c *gin.Context) {
	d := tonic.Timebox(Config.Security.Timebox)

	c.Set("user", "logged in")

	err := cookie.Drop(c, Config.Security, "user")

	d.Wait()

	if err != nil {
		c.String(http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
	}

	c.String(http.StatusOK, http.StatusText(http.StatusOK))
}

// Endpoint example, simply returns 200 if the authentication middleware
// allows access.
func Endpoint(c *gin.Context) {
	c.String(http.StatusOK, http.StatusText(http.StatusOK))
}

// Token sets the JWT on the response.
func Token(c *gin.Context) {
	jwt.Sign(c, Config.Security)
}
