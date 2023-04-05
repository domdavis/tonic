package config_test

import (
	"fmt"

	"bitbucket.org/idomdavis/gofigure"
	"bitbucket.org/idomdavis/tonic/config"
)

func ExampleTLS_Register() {
	c := gofigure.NewConfiguration("")
	t := &config.TLS{}

	t.Register(c)

	err := c.ParseUsing([]string{
		"--tls-key", "./tls.key", "--tls-certificate", "./tls.cert",
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
	// TLS settings
	//   Key Path: ./tls.key
	//   Certificate Path: ./tls.cert
	//
	// usage:
	//   Key Path [JSON key: "tls-key", env TLS-KEY, --tls-key]
	//     Path to the TLS key
	//
	//   Certificate Path [JSON key: "tls-certificate", env TLS-CERTIFICATE, --tls-certificate]
	//     Path to the TLS certificate
}

func ExampleTLS_Secure() {
	tls := &config.TLS{}

	fmt.Println(tls.Secure())

	tls.Key = "key"
	tls.Certificate = "cert"

	fmt.Println(tls.Secure())

	// Output:
	// false
	// true
}
