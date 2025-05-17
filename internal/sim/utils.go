package sim

import (
	"encoding/binary"
	"fmt"
	"go/types"
	"math"
	"strings"

	"golang.org/x/exp/constraints"
)

type Num interface {
	constraints.Float
}

func ConvertDtype(strType string) (types.Type, error) {
	strType = strings.Trim(strType, " ")
	strType = strings.ToLower(strType)
	tObj := types.Universe.Lookup(strType)
	if tObj == nil {
		return nil, fmt.Errorf("unknown type: %s", strType)
	}
	return tObj.Type(), nil
}

func WriteBits(buf []byte, offset int, val any, size int) error {
	var uval uint64
	switch v := val.(type) {
	case int8:
		uval = uint64(v)
	case int16:
		uval = uint64(v)
	case int32:
		uval = uint64(v)
	case int64:
		uval = uint64(v)
	case int:
		uval = uint64(v)
	case uint8:
		uval = uint64(v)
	case uint16:
		uval = uint64(v)
	case uint32:
		uval = uint64(v)
	case uint64:
		uval = v
	case float32:
		uval = uint64(math.Float32bits(v))
	case float64:
		uval = math.Float64bits(v)
	case bool:
		uval = 0
		if v {
			uval = 1
		}
	default:
		return fmt.Errorf("unable to write value to bits - type not allowed")
	}
	for i := 0; i < size; i++ {
		bit := (uval >> (size - i - 1)) & 1
		byteIndex := (offset + i) / 8
		bitIndex := 7 - ((offset + i) % 8)

		if bitIndex >= len(buf) {
			return fmt.Errorf("buffer overflow - %d bit not in %d size buf", bitIndex, len(buf))
		}
		if bit == 1 {
			buf[byteIndex] |= (1 << bitIndex)
		} else {
			buf[byteIndex] &^= (1 << bitIndex)
		}
	}
	return nil
}

func WriteBitsStr(buf []byte, offset int, value string, size int) error {
	for i, r := range value {
		err := WriteBits(buf, offset+i*size, r, size)
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteBitsArr(buf []byte, offset int, value []uint64, size int) error {
	for i, v := range value {
		err := WriteBits(buf, offset+i*size, v, size)
		if err != nil {
			return err
		}
	}
	return nil
}

func Htons(s *uint16) []byte {
	netBytes := make([]byte, 2)

	binary.BigEndian.PutUint16(netBytes, *s)

	return netBytes
}

func Htonl(l *uint32) []byte {
	netBytes := make([]byte, 4)

	binary.BigEndian.PutUint32(netBytes, *l)

	return netBytes
}

func Htond(d *uint64) []byte {
	netBytes := make([]byte, 8)

	binary.BigEndian.PutUint64(netBytes, *d)

	return netBytes
}