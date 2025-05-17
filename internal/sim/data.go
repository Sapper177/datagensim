package sim

import (
	"log"
	"strconv"

	"github.com/Sapper177/datagensim/pkg/engine"
)

type dtype uint8

const (
	D_BOOL    dtype = 0
	D_INT     dtype = 1
	D_INT8    dtype = 2
	D_INT16   dtype = 3
	D_INT32   dtype = 4
	D_INT64   dtype = 5
	D_UINT    dtype = 6
	D_UINT8   dtype = 7
	D_UINT16  dtype = 8
	D_UINT32  dtype = 9
	D_UINT64  dtype = 10
	D_FLOAT32 dtype = 11
	D_FLOAT64 dtype = 12
	D_STRING  dtype = 13
	D_BYTES   dtype = 14
)

func selectDtype(dtypeStr string) dtype {
	switch dtypeStr {
	case "bool":
		return D_BOOL
	case "int":
		return D_INT
	case "int8":
		return D_INT8
	case "int16":
		return D_INT16
	case "int32":
		return D_INT32
	case "int64":
		return D_INT64
	case "uint":
		return D_UINT
	case "uint8":
		return D_UINT8
	case "uint16":
		return D_UINT16
	case "uint32":
		return D_UINT32
	case "uint64":
		return D_UINT64
	case "float32":
		return D_FLOAT32
	case "float64":
		return D_FLOAT64
	case "string":
		return D_STRING
	case "bytes":
		return D_BYTES
	default:
		log.Printf("Unknown data type: %s. Defaulting to int32.", dtypeStr)
		return D_INT32
	}
}

type dataPoint interface {
	appendData(buf []byte, val any) error
	update(val any) (any, string)
}

type dataPointFloat struct {
	dtype  dtype
	eng    *engine.NumEngine64
	offset uint16
	size   uint16 // in bits
}

func newDataPointFloat(dtype dtype, numEng *engine.NumEngine64, offset uint16, size uint16) *dataPointFloat {
	return &dataPointFloat{
		dtype:  dtype,
		eng:    numEng,
		offset: offset,
		size:   size,
	}
}
func (d *dataPointFloat) appendData(buf []byte, val any) error {
	return WriteBits(buf, int(d.offset), val, int(d.size))
}
func (d *dataPointFloat) update(val any) (any, string) {
	newVal := 0.0
	switch v := val.(type) {
	case float64:
		newVal = d.eng.Update(v)
	}
	return newVal, strconv.FormatFloat(newVal, 'f', -1, 64)
}

type dataPointInt struct {
	dtype  dtype
	eng    *engine.NumEngineInt
	offset uint16
	size   uint16 // in bits
}

func newDataPoint32(dtype dtype, numEng *engine.NumEngineInt, offset uint16, size uint16) *dataPointInt {
	return &dataPointInt{
		dtype:  dtype,
		eng:    numEng,
		offset: offset,
		size:   size,
	}
}
func (d *dataPointInt) appendData(buf []byte, val any) error {
	return WriteBits(buf, int(d.offset), val, int(d.size))
}
func (d *dataPointInt) update(val any) (any, string) {
	var newVal int64 = 0
	switch v := val.(type) {
	case int64:
		newVal = d.eng.Update(v)
	}
	return newVal, strconv.FormatInt(newVal, 10)
}

type strDataPoint struct {
	dtype  dtype
	eng    *engine.StrEngine
	offset uint16
	size   uint16 // length of string
}

func newStrDataPoint(dtype dtype, strEng *engine.StrEngine, offset uint16, size uint16) *strDataPoint {
	return &strDataPoint{
		dtype:  dtype,
		eng:    strEng,
		offset: offset,
		size:   size,
	}
}
func (d *strDataPoint) appendData(buf []byte, val any) error {
	valStr := ""
	switch v := val.(type) {
	case string:
		valStr = v
	}
	return WriteBitsStr(buf, int(d.offset), valStr, int(d.size))
}
func (d *strDataPoint) update(val any) (any, string) {
	newVal := ""
	switch v := val.(type) {
	case string:
		newVal = d.eng.Update(v)
	}
	return newVal, newVal
}

type boolDataPoint struct {
	eng    *engine.BoolEngine
	offset uint16
	size   uint16 // length of string
}

func newBoolDataPoint(boolEng *engine.BoolEngine, offset uint16, size uint16) *boolDataPoint {
	return &boolDataPoint{
		eng:    boolEng,
		offset: offset,
		size:   size,
	}
}
func (d *boolDataPoint) appendData(buf []byte, val any) error {
	return WriteBits(buf, int(d.offset), val, int(d.size))
}
func (d *boolDataPoint) update(val any) (any, string) {
	newVal := false
	switch v := val.(type) {
	case bool:
		newVal = d.eng.Update(v)
	}
	s := "0"
	if newVal {
		s = "1"
	}
	return newVal, s
}
