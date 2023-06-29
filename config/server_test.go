package config_test

import (
	"fmt"

	"github.com/domdavis/gofigure"
	"github.com/domdavis/tonic/config"
)

func ExampleServer_Register() {
	c := gofigure.NewConfiguration("")

	server := &config.Server{}
	server.Register(c)

	err := c.ParseUsing([]string{})

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
	// Server settings
	//   Server Port: 8000
	//   TLS Certificate Path: UNSET
	//   TLS Kep Path: UNSET
	//   Dev Mode: false
	//   Write Timeout: 1m0s
	//   Read Timeout: 1m0s
	//
	// usage:
	//   Server Port [JSON key: "server-port", env SERVER_PORT, --server-port]
	//     HTTP Port to listen on (default: 8000)
	//
	//   Templates path [JSON key: "server-templates", env SERVER_TEMPLATES, --server-templates]
	//     Location of the server templates. Leave blank if not required
	//
	//   Static content path [JSON key: "server-static-content", env SERVER_STATIC_CONTENT, --server-static-content]
	//     Location of the static content. Leave blank if not required
	//
	//   TLS Certificate Path [JSON key: "server-cert-path", env SERVER_CERT_PATH, --server-cert-path]
	//     Path to TLS certificate file
	//
	//   TLS Kep Path [JSON key: "server-key-path", env SERVER_KEY_PATH, --server-key-path]
	//     Path to TLS key file
	//
	//   Dev Mode [--dev]
	//     Set development mode (default: false)
	//
	//   Write Timeout [JSON key: "server-write-timeout", env SERVER_WRITE_TIMEOUT, --server-write-timeout]
	//     Write timeout for the server (default: 1m0s)
	//
	//   Read Timeout [JSON key: "server-read-timeout", env SERVER_READ_TIMEOUT, --server-read-timeout]
	//     Read timeout for the server (default: 1m0s)
}

func ExampleServer_TLS() {
	s := config.Server{}

	fmt.Println(s.TLS())

	s.Key = "tls.key"
	s.Certificate = "tls.crt"

	fmt.Println(s.TLS())

	// Output:
	// false
	// true
}
