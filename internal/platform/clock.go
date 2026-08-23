package platform

import "time"

type Clock interface {
	Now() time.Time
	After(time.Duration) <-chan time.Time
}
type RealClock struct{}

func (RealClock) Now() time.Time                         { return time.Now() }
func (RealClock) After(time.Duration) <-chan time.Time { return nil }
