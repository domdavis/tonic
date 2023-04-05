package config

import (
	"fmt"
	"strings"
	"time"

	"bitbucket.org/idomdavis/gofigure"
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

// DefaultLimit for the limiter.
const DefaultLimit = int64(100)

// DefaultTTL for the limiter.
const DefaultTTL = time.Hour

// Register the limiter.
func (l *Limiter) Register(c *gofigure.Configuration) {
	var (
		separator string
		spacer    string
	)

	if l.Name != "" {
		separator = "-"
		spacer = " "
	}

	flag := strings.ToLower(l.Name)
	group := c.Group("Limiter settings")

	group.Add(gofigure.Optional(fmt.Sprintf("%s%sLimiter Name", l.Name, spacer),
		fmt.Sprintf("%s%slimiter-name", flag, separator), &l.Name, l.Name,
		gofigure.Reference, gofigure.HideValue,
		"Name of the limiter, should be set when the limiter is instantiated"))
	group.Add(gofigure.Optional(fmt.Sprintf("%s%sLimiter Limit", l.Name, spacer),
		fmt.Sprintf("%s%slimiter-limit", flag, separator), &l.Limit, DefaultLimit,
		gofigure.NamedSources, gofigure.ReportValue,
		"Number of times something needs to be seen before the limiter trips"))
	group.Add(gofigure.Optional(fmt.Sprintf("%s%sLimiter TTL", l.Name, spacer),
		fmt.Sprintf("%s%slimiter-ttl", flag, separator), &l.TTL, DefaultTTL,
		gofigure.NamedSources, gofigure.ReportValue,
		"Time period the limiter will check over before the limiter trips"))
}
