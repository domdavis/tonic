package tonic_test

import (
	"fmt"
	"time"

	"github.com/domdavis/tonic"
)

func ExampleTimebox() {
	start := time.Now()
	timebox := tonic.Timebox(time.Millisecond)
	timebox.Wait()

	fmt.Println(time.Since(start) >= time.Millisecond)

	// Output:
	// true
}
