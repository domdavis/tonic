package register_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/domdavis/tonic/middleware"
	"github.com/domdavis/tonic/register"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ExamplePing() {
	router := gin.New()

	register.Ping(router)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/ping", nil)

	if err != nil {
		fmt.Println(err)
	}

	router.ServeHTTP(w, req)
	fmt.Println(w.Code, w.Body.String())

	// Output:
	// 200 pong
}

func ExampleSilentPing() {
	logger := middleware.LogrusReporter(logrus.StandardLogger())
	router := gin.New()
	router.Use(logger.Log)

	register.SilentPing(router, logger)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/ping", nil)

	if err != nil {
		fmt.Println(err)
	}

	router.ServeHTTP(w, req)
	fmt.Println(w.Code, w.Body.String())

	// Output:
	// 200 pong
}
