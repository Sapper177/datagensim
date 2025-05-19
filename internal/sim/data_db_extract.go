package sim

import (
	"log"
	"strconv"

	"github.com/Sapper177/datagensim/pkg/config"
)

type dbExtract struct {
	offset uint16
	size   uint16 // in bits
	dtype  string
	min    string
	max    string
	step   string
	freq   float64
}

func newDBExtract(id string, dataid string, dataInfo map[string]string, info map[string]string) *dbExtract {
	// Convert offset and size to int
	offset, err := strconv.Atoi(dataInfo["offset"])
	if err != nil {
		log.Printf("Error converting offset for Payload (%s) - data ID (%s): %s", id, dataid, err)
		return nil
	}
	size, err := strconv.Atoi(dataInfo["size"])
	if err != nil {
		log.Printf("Error converting size for Payload (%s) - data ID (%s): %s", id, dataid, err)
		return nil
	}

	// Convert string type to go Type
	t, err := convertDtype(dataInfo["type"])
	if err != nil {
		log.Printf("Error converting data type for Payload (%s) - data ID (%s): %s", id, dataid, err)
		return nil
	}

	switch t.String() {
	case "bool", "string", "byte":
		// do nothing
		return &dbExtract{
			offset: uint16(offset),
			size:   uint16(size),
			dtype:  t.String(),
		}

	}
	min, ok := info["min"]
	if !ok {
		log.Printf("Error converting min for Payload (%s) - data ID (%s): %s", id, dataid, err)
		min = "0"
	}
	max, ok := info["max"]
	if !ok {
		log.Printf("Error converting max for Payload (%s) - data ID (%s): %s", id, dataid, err)
		max = "1000"
	}
	step, ok := info["step"]
	if !ok {
		log.Printf("Error converting step for Payload (%s) - data ID (%s): %s", id, dataid, err)
		step = "1"
	}
	freq, err := strconv.ParseFloat(info["frequency"], 64)
	if err != nil {
		log.Printf("Error converting frequency for Payload (%s) - data ID (%s): %s", id, dataid, err)
		freq = config.FREQ_DEFAULT
	}
	return &dbExtract{
		offset: uint16(offset),
		size:   uint16(size),
		dtype:  t.String(),
		min:    min,
		max:    max,
		step:   step,
		freq:   freq,
	}
}
