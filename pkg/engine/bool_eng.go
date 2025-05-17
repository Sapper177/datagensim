package engine

import (
	"time"
)

type BoolEngine struct {
	Frequency  time.Duration // in milliseconds
	lastUpdate time.Time
}
// NewBoolEngine creates a new BoolEnging instance with the specified frequency
func NewBoolEngine(freq time.Duration) *BoolEngine {
	return &BoolEngine{
		Frequency: freq,
	}
}

func (b *BoolEngine) Update(val bool) bool {
	b.lastUpdate = time.Now()
	// Calculate the time since the last update
	elapsed := time.Since(b.lastUpdate)
	if elapsed < b.Frequency {
		return val
	}
	return !val
}
