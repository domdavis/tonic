package config

import (
	"bitbucket.org/idomdavis/goconfigure"
)

// TLS options.
type TLS struct {
	// Key is the path to the TLS key.
	Key string `json:"Key Path"`

	// Certificate is the path to the TLS Certificate.
	Certificate string `json:"Certificate Path"`
}

// Description for the TLS settings.
func (t *TLS) Description() string {
	return "TLS settings"
}

// Register the TLS options.
func (t *TLS) Register(opts goconfigure.OptionSet) {
	opts.Add(opts.Option(&t.Key, "", "tls-key", "Path to the TLS key"))
	opts.Add(opts.Option(&t.Certificate, "", "tls-certificate",
		"Path to the TLS certificate"))
}

// Data for the TLS settings.
func (t *TLS) Data() interface{} {
	return TLS{
		Key:         goconfigure.Sanitise(t.Key, t.Key, goconfigure.UNSET),
		Certificate: goconfigure.Sanitise(t.Certificate, t.Certificate, goconfigure.UNSET),
	}
}

// Secure returns true if TLS is configured.
func (t *TLS) Secure() bool {
	return t.Key != "" && t.Certificate != ""
}
