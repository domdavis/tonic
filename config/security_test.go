package config_test

import (
	"fmt"

	"bitbucket.org/idomdavis/goconfigure"
	"bitbucket.org/idomdavis/tonic/config"
)

func ExampleSecurity_Register() {
	s := goconfigure.NewSettings("TEST")
	s.Add(&config.Security{})

	err := s.ParseUsing([]string{
		"-secret", "1234", "--session-ttl", "12h", "--login-timebox", "1s",
	}, goconfigure.ConsoleReporter{})

	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// Security settings: [Secret:SET, Session TTL:12h0m0s, Timebox:1s]
}
