package tonic

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

// Get a value from the given Context in a type safe manner. If no value is
// associated with the key then the zero value for the type will be returned. If
// the type cannot be cast to type T then Get will panic.
func Get[T any](ctx *gin.Context, key string) T {
	var t T

	if v, ok := ctx.Get(key); !ok {
		return t
	} else if t, ok = v.(T); !ok {
		panic(fmt.Sprintf("cannot get %T from context, value is %T", t, v))
	}

	return t
}

// Value returns the context value in a type safe manner. If no value is
// associated with the key then the zero value for the type will be returned.
// If the type cannot be cast to type T then Value will panic.
func Value[T any](ctx context.Context, key string) T {
	var (
		ok bool
		t  T
	)

	v := ctx.Value(key)

	if v == nil {
		return t
	} else if t, ok = v.(T); !ok {
		panic(fmt.Sprintf("cannot get %T from context, value is %T", t, v))
	}

	return t
}
