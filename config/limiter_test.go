package config_test

import (
	"fmt"

	"bitbucket.org/idomdavis/goconfigure"
	"bitbucket.org/idomdavis/tonic/config"
)

func ExampleLimiter_Register() {
	s := goconfigure.NewSettings("TEST")
	s.Add(&config.Limiter{})
	s.Add(&config.Limiter{Name: "Sub"})

	err := s.ParseUsing([]string{
		"--limiter-limit", "10", "--limiter-ttl", "1m",
		"--sub-limiter-limit", "20", "--sub-limiter-ttl", "1h",
	}, goconfigure.ConsoleReporter{})

	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// Limiter Settings: [Limit:10, Name:UNSET, TTL:1m0s]
	// Sub Limiter Settings: [Limit:20, Name:Sub, TTL:1h0m0s]
}
