package tonic

import "time"

// Deadline holds the end point for a Timebox.
type Deadline <-chan time.Time

// Timebox an action. Timebox will return a Deadline which can be waited on
// until the given duration has passed.
func Timebox(duration time.Duration) Deadline {
	return time.NewTimer(duration).C
}

// Wait until this Deadline is past.
func (d Deadline) Wait() {
	<-d
}
