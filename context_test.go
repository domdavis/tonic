package tonic_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bitbucket.org/idomdavis/tonic"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func ExampleGet() {
	var s string

	ctx := &gin.Context{}
	ctx.Set("k", "v")

	s = tonic.Get[string](ctx, "k")

	fmt.Println(s)

	// Output:
	// v
}

func ExampleValue() {
	var s string

	ctx := context.Background()
	ctx = context.WithValue(ctx, "k", "v")

	s = tonic.Value[string](ctx, "k")

	fmt.Println(s)

	// Output:
	// v
}

func ExampleConvert() {
	g := &gin.Context{Request: httptest.NewRequest(http.MethodGet, "http://example.com", nil)}
	g.Set("k", "v")

	ctx := tonic.Convert(g)

	fmt.Println(tonic.Value[string](ctx, "k"))

	// Output:
	// v
}

func TestGet(t *testing.T) {
	t.Run("Unset values with nil zero values will return nil", func(t *testing.T) {
		t.Parallel()

		p := tonic.Get[*string](&gin.Context{}, "k")

		assert.Nil(t, p)
	})

	t.Run("Unset values with non-nil zero values will return zero value", func(t *testing.T) {
		t.Parallel()

		s := tonic.Get[string](&gin.Context{}, "k")

		assert.Empty(t, s)

		i := tonic.Get[int](&gin.Context{}, "k")

		assert.Zero(t, i)
	})

	t.Run("Invalid type casting will panic", func(t *testing.T) {
		t.Parallel()

		ctx := &gin.Context{}
		ctx.Set("k", "v")

		assert.Panics(t, func() {
			tonic.Get[int](ctx, "k")
		})
	})
}

func TestValue(t *testing.T) {
	t.Run("Unset values will return zero value", func(t *testing.T) {
		t.Parallel()

		s := tonic.Value[string](context.Background(), "k")

		assert.Empty(t, s)

		i := tonic.Value[int](context.Background(), "k")

		assert.Zero(t, i)
	})

	t.Run("Invalid type casting will panic", func(t *testing.T) {
		t.Parallel()

		ctx := context.WithValue(context.Background(), "k", "v")

		assert.Panics(t, func() {
			tonic.Value[int](ctx, "k")
		})
	})
}
