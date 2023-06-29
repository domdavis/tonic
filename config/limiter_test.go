package config_test

import (
	"fmt"

	"github.com/domdavis/gofigure"
	"github.com/domdavis/tonic/config"
)

func ExampleLimiter_Register() {
	c := gofigure.NewConfiguration("")

	limiter := &config.Limiter{}
	subLimiter := &config.Limiter{Name: "Sub"}

	limiter.Register(c)
	subLimiter.Register(c)

	err := c.ParseUsing([]string{
		"--limiter-limit", "10", "--limiter-ttl", "1m",
		"--sub-limiter-limit", "20", "--sub-limiter-ttl", "1h",
	})

	if err != nil {
		fmt.Println(c.Format(err))
	}

	for _, line := range c.Report() {
		fmt.Println(line.Name)

		for k, v := range line.Values {
			fmt.Printf("  %s: %v\n", k, v)
		}
	}

	fmt.Println()
	fmt.Println(c.Usage())

	// Unordered Output:
	// Limiter settings
	//   Limiter Limit: 10
	//   Limiter TTL: 1m0s
	//   Sub Limiter Limit: 20
	//   Sub Limiter TTL: 1h0m0s
	//
	// usage:
	//   Limiter Limit [JSON key: "limiter-limit", env LIMITER-LIMIT, --limiter-limit]
	//     Number of times something needs to be seen before the limiter trips (default: 100)
	//
	//   Limiter TTL [JSON key: "limiter-ttl", env LIMITER-TTL, --limiter-ttl]
	//     Time period the limiter will check over before the limiter trips (default: 1h0m0s)
	//
	//   Sub Limiter Limit [JSON key: "sub-limiter-limit", env SUB-LIMITER-LIMIT, --sub-limiter-limit]
	//     Number of times something needs to be seen before the limiter trips (default: 100)
	//
	//   Sub Limiter TTL [JSON key: "sub-limiter-ttl", env SUB-LIMITER-TTL, --sub-limiter-ttl]
	//     Time period the limiter will check over before the limiter trips (default: 1h0m0s)
}
