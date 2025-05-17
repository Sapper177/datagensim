package engine

import (
	"math"
	"strings"
	"time"
)

type StrEngine struct {
	EngType    string
	Size       uint16
	Frequency  float64 // in milliseconds
	Phase      float64
	lastUpdate time.Time
	starIndex  int
	startTime  time.Time
}

// NewStrEngine creates a new StrEngine instance with the specified parameters
func NewStrEngine(size uint16, freq float64, phase float64) *StrEngine {
	return &StrEngine{
		Size:      size,
		Frequency: freq,
		Phase:     phase,
		startTime: time.Now(),
	}
}

// Update returns the next value depending on the type of engine
func (e *StrEngine) Update(val string) string {
	// Convert val to float32
	switch e.EngType {
	case "sin":
		return e.UpdateSin(val)
	case "ramp":
		return e.UpdateRamp(val)
	default: // aka static
		return val
	}
}

// GenerateString creates the animated string with the '*' at the calculated position.
func (e *StrEngine) GenerateString() string {
	elapsedTime := time.Since(e.startTime).Seconds() // Time elapsed in seconds

	// Calculate the sine wave value (-1 to 1)
	sineValue := math.Sin(elapsedTime*e.Frequency + e.Phase)

	// Map the sine value to an index within the string size
	position := float64(e.Size/2) + float64(e.Size)/2*sineValue

	// Round the position to get an integer index
	index := int(math.Floor(position))

	// Clamp the index to ensure it stays within the string bounds [0, StringSize - 1]
	index = int(math.Max(0, math.Min(float64(e.Size-1), float64(index))))
	e.starIndex = index

	// Create a string of '-' characters
	animatedString := strings.Repeat("-", int(e.Size))

	// Replace the character at the calculated index with '*'
	// Go strings are immutable, so we build a new string or byte slice
	runes := []rune(animatedString) // Use runes to handle potential multi-byte characters (though not needed for '-')
	runes[index] = '*'

	return string(runes)
}

func (b *StrEngine) UpdateSin(val string) string {
	b.lastUpdate = time.Now()
	// Calculate the time since the last update
	elapsed := time.Since(b.lastUpdate).Milliseconds()
	if elapsed < int64(b.Frequency) {
		return val
	}
	return b.GenerateString()
}

func (b *StrEngine) UpdateRamp(val string) string {
	b.lastUpdate = time.Now()
	// Calculate the time since the last update
	elapsed := time.Since(b.lastUpdate).Milliseconds()
	if elapsed < int64(b.Frequency) {
		return val
	}

	// Update start index
	b.starIndex++
	if b.starIndex >= int(b.Size) {
		b.starIndex = 0
	}

	// Create a string of '-' characters
	animatedString := strings.Repeat("-", int(b.Size))

	// Replace the character at the calculated index with '*'
	runes := []rune(animatedString)
	runes[b.starIndex] = '*'

	return string(runes)
}
