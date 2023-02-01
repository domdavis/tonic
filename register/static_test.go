package register_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"bitbucket.org/idomdavis/tonic"
	"bitbucket.org/idomdavis/tonic/register"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ExampleStatic() {
	router := gin.New()
	router.Use(tonic.NewLogger(logrus.StandardLogger()).Log)

	if err := register.Static("./invalid", router); err != nil {
		fmt.Println(err)
	}

	if err := register.Static("./testdata", router); err != nil {
		fmt.Println(err)
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/root", nil)

	router.ServeHTTP(w, req)
	fmt.Println(w.Code, strings.TrimSpace(w.Body.String()), err)

	w = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, "/dir/branch", nil)

	router.ServeHTTP(w, req)
	fmt.Println(w.Code, strings.TrimSpace(w.Body.String()), err)

	w = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, "/dir/root", nil)

	router.ServeHTTP(w, req)
	fmt.Println(w.Code, strings.TrimSpace(w.Body.String()), err)

	// Output:
	// failed to read ./invalid: open ./invalid: no such file or directory
	// 200 /root <nil>
	// 200 /dir/branch <nil>
	// 404  <nil>
}
