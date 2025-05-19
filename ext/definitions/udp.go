package definitions

import (
	"encoding/binary"
	"hash/crc32"
	"math"
	"math/rand" // For random byte example
	"time"
)

// --- Example dynamic Header functions ---

// GetCurrentTimestampUint32 returns the current Unix timestamp as a big-endian uint32 byte slice.
func GetCurrentTimestampUint32() ([]byte, uint16, error) {
	timestamp := uint32(time.Now().Unix()) // Get current Unix timestamp
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, timestamp) // Convert to network byte order
	return buf, 4, nil
}

// GetRandomByte returns a single random byte.
func GetRandomByte() ([]byte, uint16, error) {
    // Seed the random number generator (do this once in your application startup)
    // For simplicity here, we'll seed it based on time.
    return []byte{byte(rand.Intn(256))}, 4, nil // Get a random integer between 0 and 255
}

// example header definition with payload id
func NewUdpHeader(id uint32) *Header {
	idConst := Uint32Constant(id)
	return &Header{
		Elements: []ByteSource{
			ByteConstant(0xAA),                 // A fixed start byte
			StringConstant("VERSION_1"),        // A fixed version string
			Uint16Constant(0x1234),             // A fixed 16-bit identifier (big-endian)
			FuncSource{Fn: GetCurrentTimestampUint32}, // Dynamic timestamp
			ByteConstant(0xBB),                 // Another fixed byte
			Uint32Constant(idConst),			// Add payload id
            FuncSource{Fn: GetRandomByte},      // Dynamic random byte
            Uint32Constant(0x56789ABC),         // A fixed 32-bit value (big-endian)
		},
	}
}


// example Footer definition with CRC 32 bit checksum
type CRC32 struct {payload []byte}

func (c *CRC32) CalculateCRC32() uint32 {
	table := crc32.MakeTable(crc32.IEEE)
	return crc32.Checksum(c.payload, table)
}

func (c *CRC32) Bytes() ([]byte, uint16, error) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, c.CalculateCRC32()) // Use network byte order (big-endian)
	return buf, 4, nil
}

func NewUdpFooter(payload []byte) *Footer {
	return &Footer{
		Elements: []ByteSource{
			&CRC32{payload: payload},
		},
	}
}
