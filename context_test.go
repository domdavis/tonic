package tonic_test

import (
	"context"
	"fmt"
	"testing"

	"bitbucket.org/idomdavis/tonic"
	"github.com/stretchr/testify/assert"
)

func ExampleGet() {
	var s string

	ctx := context.WithValue(context.Background(), "k", "v")

	s = tonic.Get[string](ctx, "k")

	fmt.Println(s)

	// Output:
	// v
}

func TestGet(t *testing.T) {
	t.Run("Unset values with nil zero values will return nil", func(t *testing.T) {
		t.Parallel()

		p := tonic.Get[*string](context.Background(), "k")

		assert.Nil(t, p)
	})

	t.Run("Unset values with non-nil zero values will return zero value", func(t *testing.T) {
		t.Parallel()

		s := tonic.Get[string](context.Background(), "k")

		assert.Empty(t, s)

		i := tonic.Get[int](context.Background(), "k")

		assert.Zero(t, i)
	})

	t.Run("Invalid type casting will panic", func(t *testing.T) {
		t.Parallel()

		ctx := context.WithValue(context.Background(), "k", "v")

		assert.Panics(t, func() {
			tonic.Get[int](ctx, "k")
		})
	})
}
