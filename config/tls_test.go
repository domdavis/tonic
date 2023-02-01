package config_test

import (
	"fmt"

	"bitbucket.org/idomdavis/goconfigure"
	"bitbucket.org/idomdavis/tonic/config"
)

func ExampleTLS_Register() {
	s := goconfigure.NewSettings("TEST")
	s.Add(&config.TLS{})

	err := s.ParseUsing([]string{
		"-tls-key", "./tls.key", "--tls-certificate", "./tls.cert",
	}, goconfigure.ConsoleReporter{})

	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// TLS settings: [Certificate Path:./tls.cert, Key Path:./tls.key]
}
