package tonic

import "time"

// Deadline holds the end point for a Timebox.
type Deadline time.Time

// Timebox an action. Timebox will return a Deadline which can be waited on
// until the given duration has passed.
func Timebox(duration time.Duration) Deadline {
	return Deadline(time.Now().Add(duration))
}

// Wait until this Deadline is past.
func (d Deadline) Wait() {
	remaining := time.Until(time.Time(d))

	if remaining <= 0 {
		return
	}

	time.Sleep(remaining)
}
