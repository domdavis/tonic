package register_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/domdavis/tonic/register"
	"github.com/gin-gonic/gin"
)

func ExampleCoffee() {
	router := gin.New()

	register.Coffee(router)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/coffee", nil)

	if err != nil {
		fmt.Println(err)
	}

	router.ServeHTTP(w, req)
	fmt.Println(w.Code, w.Body.String())

	// Output:
	// 418 I'm a teapot
}
