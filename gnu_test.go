package tonic_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/domdavis/tonic"
	"github.com/domdavis/tonic/middleware"
	"github.com/domdavis/tonic/register"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ExampleGNU() {
	router := gin.New()
	router.Use(middleware.LogrusReporter(logrus.StandardLogger()).Log)
	router.Use(tonic.GNU)

	register.Ping(router)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/ping", nil)

	if err != nil {
		fmt.Println(err)
	}

	router.ServeHTTP(w, req)

	fmt.Println(w.Result().Header.Values(tonic.Clacks))

	// Output:
	// [GNU Terry Pratchett, Russel Winder]
}
