package engine

import (
	"math";
	"reflect";
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

func (s *SinEngine) update(val interface{}) interface{} {
	switch val.(type) {
	case float32, float64:
		s.counter += 1
		val = s.Min + (s.Max-s.Min)/2 + (s.Max-s.Min)/2*math.Sin(2*math.Pi*s.Frequency*s.counter+s.Phase)
	case int, int16, int32, int64:
		s.counter += 1
		val = int(s.Min + (s.Max-s.Min)/2 + (s.Max-s.Min)/2*math.Sin(2*math.Pi*s.Frequency*s.counter+s.Phase))
	case uint, uint16, uint32, uint64:
		s.counter += 1
		val = uint(s.Min + (s.Max-s.Min)/2 + (s.Max-s.Min)/2*math.Sin(2*math.Pi*s.Frequency*s.counter+s.Phase))
	}
	return val
}

type RampEngine struct {
    Min       float64
    Max       float64
    Step      float64
    Frequency float64
    Phase     float64
    counter   float64
}

func (s *RampEngine) update(val interface{}) interface{} {

	switch v := val.(type) {
	case float64:
		v += s.Step
		if v > s.Max {
			v = s.Min
		}
	case int, int16, int32, int64:
		var value float64 = float64(val)
		v = v. int(s.Step)
		if v > s.Max {
			v = s.Min
		}
	case uint, uint16, uint32, uint64:
		s.counter += 1
		val = uint(s.Min + (s.Max-s.Min)/2 + (s.Max-s.Min)/2*math.Sin(2*math.Pi*s.Frequency*s.counter+s.Phase))
	}
	return v
}