package register_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bitbucket.org/idomdavis/tonic/register"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func ExampleStatic() {
	router := gin.New()

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

func TestStatic(t *testing.T) {
	t.Run("Blank content roots are ignored", func(t *testing.T) {
		t.Parallel()

		router := gin.New()
		err := register.Static("", router)

		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/root", nil)

		router.ServeHTTP(w, req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
