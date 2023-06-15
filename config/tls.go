package config

import (
	"bitbucket.org/idomdavis/gofigure"
)

// TLS options.
type TLS struct {
	// Key is the path to the TLS key.
	Key string `json:"Key Path"`

	// Certificate is the path to the TLS Certificate.
	Certificate string `json:"Certificate Path"`
}

// Register the TLS options.
func (t *TLS) Register(c *gofigure.Configuration) {
	group := c.Group("TLS settings")

	group.Add(gofigure.Optional("Key Path", "tls-key", &t.Key, "",
		gofigure.NamedSources, gofigure.HideUnset, "Path to the TLS key"))
	group.Add(gofigure.Optional("Certificate Path", "tls-certificate",
		&t.Certificate, "", gofigure.NamedSources, gofigure.HideUnset,
		"Path to the TLS certificate"))
}

// Secure returns true if TLS is configured.
func (t *TLS) Secure() bool {
	return t.Key != "" && t.Certificate != ""
}
