package definitions

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// ByteSource is an interface that any header element must implement.
// It defines a method to get the byte representation of the element.
type ByteSource interface {
	Bytes() ([]byte, error)
}

// --- Implementations for Constant Values ---

// ByteConstant wraps a single byte constant.
type ByteConstant byte

// Bytes returns the byte representation of the ByteConstant.
func (b ByteConstant) Bytes() ([]byte, error) {
	return []byte{byte(b)}, nil
}

// Uint16Constant wraps a uint16 constant.
type Uint16Constant uint16

// Bytes returns the big-endian byte representation of the Uint16Constant.
func (u Uint16Constant) Bytes() ([]byte, error) {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(u)) // Use network byte order (big-endian)
	return buf, nil
}

// Uint32Constant wraps a uint32 constant.
type Uint32Constant uint32

// Bytes returns the big-endian byte representation of the Uint32Constant.
func (u Uint32Constant) Bytes() ([]byte, error) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(u)) // Use network byte order (big-endian)
	return buf, nil
}

// StringConstant wraps a string constant.
// Note: This assumes the string bytes themselves are the desired payload representation.
type StringConstant string

// Bytes returns the byte representation of the StringConstant (UTF-8 bytes).
func (s StringConstant) Bytes() ([]byte, error) {
	return []byte(s), nil
}


// --- Implementation for Function Results ---

// ByteFunc is a function type that returns bytes and an error.
type ByteFunc func() ([]byte, error)

// FuncSource wraps a ByteFunc.
type FuncSource struct {
	Fn ByteFunc
}

// Bytes calls the wrapped function and returns its result.
func (fs FuncSource) Bytes() ([]byte, error) {
	if fs.Fn == nil {
		return nil, errors.New("FuncSource has a nil function")
	}
	return fs.Fn()
}

// --- Header Structure ---

// Header defines the structure of the payload header as a sequence of ByteSource elements.
type Header struct {
	Elements []ByteSource
}

// --- Payload Builder Function ---

// BuildPayload takes a Header and serializes its elements into a byte slice.
func BuildPayload(header *Header) ([]byte, error) {
	if header == nil {
		return nil, errors.New("cannot build payload from nil header")
	}

	var payload []byte // Start with an empty byte slice

	// Iterate through each element in the header
	for i, element := range header.Elements {
		if element == nil {
             return nil, fmt.Errorf("header element at index %d is nil", i)
        }
		// Call the Bytes() method on the element to get its byte representation
		elementBytes, err := element.Bytes()
		if err != nil {
			// Return an informative error if getting bytes fails for any element
			return nil, fmt.Errorf("failed to get bytes for header element at index %d: %w", i, err)
		}
		// Append the obtained bytes to the payload
		payload = append(payload, elementBytes...)
	}

	return payload, nil
}

