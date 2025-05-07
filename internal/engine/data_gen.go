package engine

import (
	"log"
	"math"
	"time"
)

type EngineType int

const (
    Static EngineType = iota
    Ramp
    Sinusoidal
)

type SinEngine struct {
    Min       float64
    Max       float64
    Step      float64
    Frequency float64
    Phase     float64
    counter   float64
}
func (s *SinEngine) sinUpdate() float64 {
	return s.Min + (s.Max-s.Min)/2 + (s.Max-s.Min)/2*math.Sin(2*math.Pi*s.Frequency*s.counter+s.Phase)
}

func (s *SinEngine) Update(val interface{}) interface{} {
	switch val.(type) {
	case float32, float64:
		newVal := s.sinUpdate()
		return newVal

	// --- int/uint cases will lose fraction part if any but for simulation this is ok ---
	case int64:
		newVal := int64(s.sinUpdate())
		return newVal
	case int32:
		newVal := int32(s.sinUpdate())
		return newVal
	case int:
		newVal := int(s.sinUpdate())
		return newVal
	case int16:
		newVal := int16(s.sinUpdate())
		return newVal
	case uint64:
		newVal := uint64(s.sinUpdate())
		return newVal
	case uint32:
		newVal := uint32(s.sinUpdate())
		return newVal
	case uint:
		newVal := uint(s.sinUpdate())
		return newVal
	case uint16:
		newVal := uint16(s.sinUpdate())
		return newVal
	default:
		log.Printf("ERROR: Unsupported type %T for SinEngine", val)
	}
	return val
}

type RampEngine struct {
    Min       float64
    Max       float64
    Step      float64
    Frequency float64 // in milliseconds
	lastUpdate time.Time
    counter   int64
}

func (s *RampEngine) Update(val interface{}) interface{} {
	s.counter++
	s.lastUpdate = time.Now()

	switch v := val.(type) {
	case float64:
		newVal := v + s.Step
		if newVal > s.Max {
			newVal = s.Min
		}
		return newVal
	case float32:
		newVal := v + float32(s.Step)
		if newVal > float32(s.Max) {
			newVal = float32(s.Min)
		}
		return newVal

	// --- int/uint cases will lose fraction part if any but for simulation this is ok ---
	case int64:
		newVal := v + int64(s.Step)
		if newVal > int64(s.Max) {
			newVal = int64(s.Min)
		}
		return newVal
	case int32:
		newVal := v + int32(s.Step)
		if newVal > int32(s.Max) {
			newVal = int32(s.Min)
		}
		return newVal
	case int:
		newVal := v + int(s.Step)
		if newVal > int(s.Max) {
			newVal = int(s.Min)
		}
		return newVal
	case int16:
		newVal := v + int16(s.Step)
		if newVal > int16(s.Max) {
			newVal = int16(s.Min)
		}
		return newVal
	case uint64:
		newVal := v + uint64(s.Step)
		if newVal > uint64(s.Max) {
			newVal = uint64(s.Min)
		}
		return newVal
	case uint32:
		newVal := v + uint32(s.Step)
		if newVal > uint32(s.Max) {
			newVal = uint32(s.Min)
		}
		return newVal
	case uint:
		newVal := v + uint(s.Step)
		if newVal > uint(s.Max) {
			newVal = uint(s.Min)
		}
		return newVal
	case uint16:
		newVal := v + uint16(s.Step)
		if newVal > uint16(s.Max) {
			newVal = uint16(s.Min)
		}
		return newVal
	default:
		log.Printf("ERROR: Unsupported type %T for RampEngine", val)
	}
	return val
}