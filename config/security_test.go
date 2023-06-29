package config_test

import (
	"fmt"

	"github.com/domdavis/gofigure"
	"github.com/domdavis/tonic/config"
)

func ExampleSecurity_Register() {
	c := gofigure.NewConfiguration("")
	s := &config.Security{}

	s.Register(c)

	err := c.ParseUsing([]string{
		"--secret", "1234", "--session-ttl", "12h", "--login-timebox", "1s",
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
	// Security settings
	//   Cookie Domain: UNSET
	//   Session TTL: 12h0m0s
	//   Login Timebox: 1s
	//   JWT Secret: SET
	//
	// usage:
	//   JWT Secret [JSON key: "secret", env SECRET, --secret]
	//     Secret used to encrypt tokens. The default is to use a random secret. (default: <random>)
	//
	//   Cookie Domain [JSON key: "domain", env DOMAIN, --domain]
	//     Cookie domain, leave blank to allow insecure cookies
	//
	//   Session TTL [JSON key: "session-ttl", env SESSION-TTL, --session-ttl]
	//     TTL for sessions (default: 12h0m0s)
	//
	//   Login Timebox [JSON key: "login-timebox", env LOGIN-TIMEBOX, --login-timebox]
	//     Minimum time it will take for a login attempt to return (default: 500ms)
}

func ExampleSecurity_Secure() {
	s := config.Security{}

	fmt.Println(s.Secure())

	s.Domain = "localhost"

	fmt.Println(s.Secure())

	s.Domain = "example.com"

	fmt.Println(s.Secure())

	// Output:
	// false
	// false
	// true
}
