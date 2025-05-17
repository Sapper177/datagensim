package definitions

import (
	"encoding/binary"
	"time"
    "math/rand" // For random byte example
)

// --- Example dynamic Header functions ---

// GetCurrentTimestampUint32 returns the current Unix timestamp as a big-endian uint32 byte slice.
func GetCurrentTimestampUint32() ([]byte, error) {
	timestamp := uint32(time.Now().Unix()) // Get current Unix timestamp
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, timestamp) // Convert to network byte order
	return buf, nil
}

// GetRandomByte returns a single random byte.
func GetRandomByte() ([]byte, error) {
    // Seed the random number generator (do this once in your application startup)
    // For simplicity here, we'll seed it based on time.
    return []byte{byte(rand.Intn(256))}, nil // Get a random integer between 0 and 255
}


// example header definition
func NewUdpHeader() *Header {
	return &Header{
		Elements: []ByteSource{
			ByteConstant(0xAA),                 // A fixed start byte
			StringConstant("VERSION_1"),        // A fixed version string
			Uint16Constant(0x1234),             // A fixed 16-bit identifier (big-endian)
			FuncSource{Fn: GetCurrentTimestampUint32}, // Dynamic timestamp
			ByteConstant(0xBB),                 // Another fixed byte
            FuncSource{Fn: GetRandomByte},      // Dynamic random byte
            Uint32Constant(0x56789ABC),         // A fixed 32-bit value (big-endian)
		},
	}
}