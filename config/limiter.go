package config

import (
	"fmt"
	"strings"
	"time"

	"bitbucket.org/idomdavis/goconfigure"
)

// Limiter configuration.
type Limiter struct {
	// Name of the limiter. This allows limiters to be namespaced if there is
	// more than one limiter being configured. Name can be blank if only one
	// limiter is being used.
	Name string

	// Limit at which the limiter will trip (i.e. the number of time something
	// needs to be seen before the limiter trips).
	Limit int64

	// TTL for the items in the limiter.
	TTL time.Duration
}

// DefaultLimit for the limiter. This is per TTL.
const DefaultLimit = int64(100)

// Description for the Limiter.
func (l *Limiter) Description() string {
	const description = "Limiter Settings"

	if l.Name == "" {
		return description
	}

	return fmt.Sprintf("%s %s", l.Name, description)
}

// Register the limiter.
func (l *Limiter) Register(opts goconfigure.OptionSet) {
	var name string

	if l.Name != "" {
		name = fmt.Sprintf("%s-", strings.ToLower(l.Name))
	}

	opts.Add(opts.Option(&l.Limit, DefaultLimit, fmt.Sprintf("%slimiter-limit", name),
		"Number of times something needs to be seen before the limiter trips"))
	opts.Add(opts.Option(&l.TTL, time.Second, fmt.Sprintf("%slimiter-ttl", name),
		"Time-to-live for items in the limiter"))
}

// Data for the Limiter.
func (l *Limiter) Data() any {
	return struct {
		Name  string
		Limit int64
		TTL   string
	}{
		Name:  goconfigure.Sanitise(l.Name, l.Name, goconfigure.UNSET),
		Limit: l.Limit,
		TTL:   l.TTL.String(),
	}
}
