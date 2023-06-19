package config

import (
	"time"

	"bitbucket.org/idomdavis/gofigure"
)

// Server configuration options.
type Server struct {
	// Port is the port number to run the server on. The default value is 8000.
	Port int

	// Templates path for the server. Leave blank if templates are not required.
	Templates string

	// Static content for the server. Leave blank if static content is not
	// required.
	Static string

	// Certificate path for TLS.
	Certificate string

	// Key path for TLS.
	Key string

	// Development flag, can be passed to Gin to control how it is set up.
	Development bool

	// WriteTimeout for the server. Default is 1 minute.
	WriteTimeout time.Duration

	// ReadTimeout for the server. Default is 1 minute.
	ReadTimeout time.Duration
}

const (
	defaultPort    = 8000
	defaultTimeout = 60 * time.Second
)

// Register the Server options.
func (s *Server) Register(c *gofigure.Configuration) {
	group := c.Group("Server settings")

	group.Add(gofigure.Optional("Server Port", "server-port", &s.Port, defaultPort,
		gofigure.NamedSources, gofigure.ReportValue, "HTTP Port to listen on"))
	group.Add(gofigure.Optional("Templates path", "server-templates",
		&s.Templates, "", gofigure.NamedSources, gofigure.HideUnset,
		"Location of the server templates. Leave blank if not required"))
	group.Add(gofigure.Optional("Static content path", "server-static-content",
		&s.Static, "", gofigure.NamedSources, gofigure.HideUnset,
		"Location of the static content. Leave blank if not required"))
	group.Add(gofigure.Optional("TLS Certificate Path", "server-cert-path",
		&s.Certificate, "", gofigure.NamedSources, gofigure.MaskUnset,
		"Path to TLS certificate file"))
	group.Add(gofigure.Optional("TLS Kep Path", "server-key-path",
		&s.Key, "", gofigure.NamedSources, gofigure.MaskUnset,
		"Path to TLS key file"))
	group.Add(gofigure.Optional("Dev Mode", "dev", &s.Development, false,
		gofigure.Flag, gofigure.ReportValue, "Set development mode"))
	group.Add(gofigure.Optional("Write Timeout", "server-write-timeout",
		&s.WriteTimeout, defaultTimeout, gofigure.NamedSources, gofigure.ReportValue,
		"Write timeout for the server"))
	group.Add(gofigure.Optional("Read Timeout", "server-read-timeout",
		&s.ReadTimeout, defaultTimeout, gofigure.NamedSources, gofigure.ReportValue,
		"Read timeout for the server"))
}

// TLS returns true if the server is configured for TLS.
func (s *Server) TLS() bool {
	return s.Certificate != "" && s.Key != ""
}
