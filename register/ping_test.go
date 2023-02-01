package register_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"bitbucket.org/idomdavis/tonic"
	"bitbucket.org/idomdavis/tonic/register"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ExamplePing() {
	router := gin.New()
	router.Use(tonic.NewLogger(logrus.StandardLogger()).Log)

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
