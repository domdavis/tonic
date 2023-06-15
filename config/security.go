package config

import (
	"time"

	"bitbucket.org/idomdavis/gofigure"
)

// Security settings.
type Security struct {
	// Secret used to encrypt JWT tokens.
	Secret string

	// Domain this service is running on.
	Domain string

	// SessionTTL is the length of time a session will be valid for before
	// it requires re-authentication.
	SessionTTL time.Duration

	// Timebox is the minimum time it will take for a login attempt to
	// return.
	Timebox time.Duration
}

// RandomSecret is used to tell tonic to use a random secret. The secret will
// be generated on startup.
const RandomSecret = "<random>"

const (
	defaultSessionTTL = time.Hour * 12
	defaultTimebox    = time.Millisecond * 500
)

// Register the Security options.
func (s *Security) Register(c *gofigure.Configuration) {
	group := c.Group("Security settings")

	group.Add(gofigure.Optional("JWT Secret", "secret",
		&s.Secret, RandomSecret, gofigure.NamedSources, gofigure.MaskSet,
		"Secret used to encrypt tokens. The default is to use a random secret."))
	group.Add(gofigure.Optional("Cookie Domain", "domain", &s.Domain, "",
		gofigure.NamedSources, gofigure.MaskUnset,
		"Cookie domain, leave blank to allow insecure cookies"))
	group.Add(gofigure.Optional("Session TTL", "session-ttl", &s.SessionTTL,
		defaultSessionTTL, gofigure.NamedSources, gofigure.ReportValue,
		"TTL for sessions"))
	group.Add(gofigure.Optional("Login Timebox", "login-timebox", &s.Timebox,
		defaultTimebox, gofigure.NamedSources, gofigure.ReportValue,
		"Minimum time it will take for a login attempt to return"))
}

// Secure returns true if a Domain is set and isn't localhost.
func (s *Security) Secure() bool {
	return s.Domain != "" && s.Domain != "localhost"
}
