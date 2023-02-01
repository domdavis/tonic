package config

import (
	"time"

	"bitbucket.org/idomdavis/goconfigure"
)

// Security settings.
type Security struct {
	// Secret used to encrypt JWT tokens.
	Secret string

	// SessionTTL is the length of time a session will be valid for before
	// it requires re-authentication.
	SessionTTL time.Duration

	// Timebox is the minimum time it will take for a login attempt to
	// return.
	Timebox time.Duration
}

// Description for the Security settings.
func (s *Security) Description() string {
	return "Security settings"
}

// Register the Security options.
func (s *Security) Register(opts goconfigure.OptionSet) {
	opts.Add(opts.Option(&s.Secret, "", "secret",
		"Secret used to encrypt tokens. Leave blank for a random secret"))
	opts.Add(opts.Option(&s.SessionTTL, time.Hour*12, "session-ttl",
		"TTL for sessions"))
	opts.Add(opts.Option(&s.Timebox, time.Millisecond*500, "login-timebox",
		"Minimum time it will take for a login attempt to return"))
}

// Data for the Security settings.
func (s *Security) Data() any {
	return struct {
		Secret     string
		SessionTTL string `json:"Session TTL"`
		Timebox    string
	}{
		Secret:     goconfigure.Sanitise(s.Secret, goconfigure.SET, goconfigure.UNSET),
		SessionTTL: s.SessionTTL.String(),
		Timebox:    s.Timebox.String(),
	}
}
