package tonic

import (
	"context"
	"fmt"
)

// Get a value from the given Context in a type safe manner. If no value is
// associated with the key then the zero value for the type will be returned. If
// the type cannot be cast to type T then Get will panic.
func Get[T any](ctx context.Context, key string) T {
	v := ctx.Value(key)

	if v == nil {
		return *new(T)
	}

	t, ok := v.(T)

	if !ok {
		panic(fmt.Errorf("cannot get %T from context, value is %T", t, v))
	}

	return t
}
