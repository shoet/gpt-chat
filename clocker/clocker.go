package clocker

import "time"

type Clocker interface {
	Now() time.Time
}

type RealClocker struct{}

func (c *RealClocker) Now() time.Time {
	return time.Now()
}

type FixedClocker struct{}

func (c *FixedClocker) Now() time.Time {
	return time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
}
