package engine

import (
	"math"
	"time"

	"golang.org/x/exp/constraints"
)

type engData[T constraints.Float] struct {
	Min        T
	Max        T
	Step       T
	Hz         T
	Frequency  time.Duration // in milliseconds
	Phase      T
	lastUpdate time.Time
}

func (s *engData[T]) getSin() T {
	s.lastUpdate = time.Now()

	// midpoint
	mid := s.Min + ((s.Max - s.Min) / 2)

	// amplitude
	amp := (s.Max - s.Min) / 2

	// 2*PI*hz*t + phase
	deg := 2*math.Pi*s.Hz*T(s.lastUpdate.Second()) + s.Phase

	// midpoint + amplitude * sin(2*PI*hz*t + phase)
	return mid + amp*T(math.Sin(float64(deg)))
}

func (s *engData[T]) getRamp(val T) T {
	s.lastUpdate = time.Now()

	// Calculate the time since the last update
	elapsed := time.Since(s.lastUpdate)
	if elapsed < s.Frequency {
		return s.Min
	}

	newVal := T(val) + s.Step
	if newVal > s.Max {
		newVal = s.Min
	}
	return newVal
}

type NumEngine64 struct {
	engType string
	engData[float64]
}

// NewNumEng64 creates a new NumEngine64 instance with the specified parameters
func NewNumEng64(min float64, max float64, step float64, freq time.Duration, phase float64, engType string) *NumEngine64 {
	hz := 1 / freq.Seconds()
	return &NumEngine64{
		engType: engType,
		engData: engData[float64]{
			Min:       min,
			Max:       max,
			Step:      step,
			Frequency: freq,
			Phase:     phase,
			Hz:        hz,
		},
	}
}

// Update returns the next value depending on the type of engine
func (e *NumEngine64) Update(val float64) float64 {
	switch e.engType {
	case "sin":
		return e.getSin()
	case "ramp":
		return e.getRamp(val)
	default:	// aka static
		return val
	}
}

type NumEngineInt struct {
	engType string
	engData[float32]
}

// NewNumEngInt creates a new NumEngineInt instance with the specified parameters
func NewNumEngInt(min float32, max float32, step float32, freq time.Duration, phase float32, engType string) *NumEngineInt {
	hz := float32(1 / freq.Seconds())
	return &NumEngineInt{
		engType: engType,
		engData: engData[float32]{
			Min:       min,
			Max:       max,
			Step:      step,
			Frequency: freq,
			Phase:     phase,
			Hz:        hz,
		},
	}
}

// Update returns the next value depending on the type of engine
func (e *NumEngineInt) Update(val int64) int64 {
	// Convert val to float32
	valF := float32(val)
	switch e.engType {
	case "sin":
		return int64(e.getSin())
	case "ramp":
		return int64(e.getRamp(valF))
	default:	// aka static
		return val
	}
}
