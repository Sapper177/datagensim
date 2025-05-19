package definitions

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

// ByteSource is an interface that any header element must implement.
// It defines a method to get the byte representation of the element.
type ByteSource interface {
	Bytes() ([]byte, uint16, error)
}

// --- Implementations for Constant Values ---

// ByteConstant wraps a single byte constant.
type ByteConstant byte

// Bytes returns the byte representation of the ByteConstant.
func (b ByteConstant) Bytes() ([]byte, uint16, error) {
	s := uint16(1)
	return []byte{byte(b)}, s, nil
}

// Uint16Constant wraps a uint16 constant.
type Uint16Constant uint16

// Bytes returns the big-endian byte representation of the Uint16Constant.
func (u Uint16Constant) Bytes() ([]byte, uint16, error) {
	s := uint16(2)
	buf := make([]byte, s)
	binary.BigEndian.PutUint16(buf, uint16(u)) // Use network byte order (big-endian)
	return buf, s, nil
}

// Uint32Constant wraps a uint32 constant.
type Uint32Constant uint32

// Bytes returns the big-endian byte representation of the Uint32Constant.
func (u Uint32Constant) Bytes() ([]byte, uint16, error) {
	s := uint16(4)
	buf := make([]byte, s)
	binary.BigEndian.PutUint32(buf, uint32(u)) // Use network byte order (big-endian)
	return buf, s, nil
}

// StringConstant wraps a string constant.
// Note: This assumes the string bytes themselves are the desired payload representation.
type StringConstant string

// Bytes returns the byte representation of the StringConstant (UTF-8 bytes).
func (s StringConstant) Bytes() ([]byte, uint16, error) {
	l := len(s)
	if l > math.MaxUint16 {
		return []byte(s), math.MaxUint16, fmt.Errorf("string larger than uint16 max")
	}
	return []byte(s), uint16(l), nil
}


// --- Implementation for Function Results ---

// ByteFunc is a function type that returns bytes and an error.
type ByteFunc func() ([]byte, uint16, error)

// FuncSource wraps a ByteFunc.
type FuncSource struct {
	Fn ByteFunc
}

// Bytes calls the wrapped function and returns its result.
func (fs FuncSource) Bytes() ([]byte, uint16, error) {
	if fs.Fn == nil {
		return nil, 0, errors.New("FuncSource has a nil function")
	}
	return fs.Fn()
}

// --- Header Structure ---

// Header defines the structure of the payload header as a sequence of ByteSource elements.
type Header struct {
	Elements []ByteSource
}

type Footer struct {
	Elements []ByteSource
}
