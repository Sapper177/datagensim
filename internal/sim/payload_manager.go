package sim

import (
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"strconv"
	"time"

	"github.com/Sapper177/datagensim/ext/definitions"
	"github.com/Sapper177/datagensim/pkg/config"
	"github.com/Sapper177/datagensim/pkg/database"
	"github.com/Sapper177/datagensim/pkg/engine"

	"github.com/google/gopacket/layers"
)

type pktType uint8

const (
	UDP pktType = iota
	TCP
	// add more
)

type payloadManager struct {
	// packet info
	pktType pktType // UDP, TCP, etc
	src     net.IP
	dst     net.IP
	srcPort layers.UDPPort
	dstPort layers.UDPPort
	srcMAC  net.HardwareAddr

	// payload info
	id       uint // payload id
	freq     time.Duration
	dpMap    map[string]dataPoint // data id -> data point
	size     uint16               //payload size
	pBuf     []byte
	lastProc time.Time
	cs       *PayloadChans
	info     *packetInfo

	// Header info
	header *definitions.Header
	hsize  uint16
	hBuf   []byte

	// Footer info
	footer *definitions.Footer
	fsize  uint16
	fBuf   []byte

	// combined full payload
	payload []byte
}

func newPayloadManager(cfg *config.Config, id string, fs time.Duration, db *database.RedisClient) *payloadManager {
	//----- Generate the payload data points -----
	// Get the list of data ids from the database
	dataids, err := db.GetPayloadData(id)
	if err != nil {
		log.Fatalf("Error getting payload data ids for ID (%s): %s", id, err)
	}
	// Create a map to hold the data points
	dps := make(map[string]dataPoint, len(dataids))

	for i := range dataids {
		// Get the data point info from the database
		dataInfo, err := db.GetData(dataids[i])
		if err != nil {
			log.Fatalf("Error getting data point info for ID (%s): %s", id, err)
		}
		// Get extra info about data point from the database
		dInfo, err := db.GetDataInfo(dataids[i])
		if err != nil {
			log.Fatalf("Error getting data point info for ID (%s): %s", id, err)
		}
		// Extract data from Database
		dbEx := newDBExtract(id, dataids[i], dataInfo, dInfo)

		dtype := selectDtype(dbEx.dtype)

		switch dtype {
		case 0:
			datapoint := newBoolDataPoint(
				engine.NewBoolEngine(time.Duration(dbEx.freq)),
				dbEx.offset,
				dbEx.size,
			)
			dps[dataids[i]] = datapoint

		case 1, 2, 3, 4, 5, 6, 7, 8, 9, 10:
			// Translate the data point info to the correct type
			min, err := strconv.ParseInt(dbEx.min, 0, 64)
			if err != nil {
				log.Printf("Error converting min for Payload (%s) - data ID (%s): %s", id, dataids[i], err)
			}
			max, err := strconv.ParseInt(dbEx.max, 0, 64)
			if err != nil {
				log.Printf("Error converting max for Payload (%s) - data ID (%s): %s", id, dataids[i], err)
			}
			step, err := strconv.ParseInt(dbEx.step, 0, 64)
			if err != nil {
				log.Printf("Error converting step for Payload (%s) - data ID (%s): %s", id, dataids[i], err)
			}

			eng := engine.NewNumEngInt(
				float32(min),
				float32(max),
				float32(step),
				time.Duration(dbEx.freq),
				float32(config.PHASE_DEFAULT),
				"sin",
			)

			datapoint := newDataPoint32(
				dtype,
				eng,
				dbEx.offset,
				dbEx.size,
			)
			dps[dataids[i]] = datapoint

		case 11, 12:
			// Translate the data point info to the correct type
			min, err := strconv.ParseFloat(dbEx.min, 64)
			if err != nil {
				log.Printf("Error converting min for Payload (%s) - data ID (%s): %s", id, dataids[i], err)
			}
			max, err := strconv.ParseFloat(dbEx.max, 64)
			if err != nil {
				log.Printf("Error converting max for Payload (%s) - data ID (%s): %s", id, dataids[i], err)
			}
			step, err := strconv.ParseFloat(dbEx.step, 64)
			if err != nil {
				log.Printf("Error converting step for Payload (%s) - data ID (%s): %s", id, dataids[i], err)
			}

			eng := engine.NewNumEng64(
				min,
				max,
				step,
				time.Duration(dbEx.freq),
				config.PHASE_DEFAULT,
				"sin",
			)

			datapoint := newDataPointFloat(
				dtype,
				eng,
				dbEx.offset,
				dbEx.size,
			)
			dps[dataids[i]] = datapoint

		case 13, 14:
			eng := engine.NewStrEngine(dbEx.size, dbEx.freq, config.PHASE_DEFAULT)
			datapoint := newStrDataPoint(dtype, eng, dbEx.offset, dbEx.size)
			dps[dataids[i]] = datapoint

		default:
			log.Printf("Unknown data type for Payload (%s) - data ID (%s): %s", id, dataids[i], dbEx.dtype)
		}
	}
	payId, err := strconv.ParseUint(id, 0, 32)
	if err != nil {
		log.Printf("Error converting id for Payload (%s): %s", id, err)
	}

	// get payload data size
	size := calcPayloadSize(dps)

	payloadBuf := make([]byte, size)

	// initialize header variables
	header := definitions.NewUdpHeader(uint32(payId), size )

	hSize, err := calcHeaderSize(*header)
	if err != nil {
		log.Printf("Error calculating error size: %s", err)
	}
	hBuf := make([]byte, hSize)

	// create tmp payload buffer for footer
	totalPayloadSize := size + hSize

	
	// footer := definitions.NewUdpFooter(uint32(payId))

	// hSize, err := calcHeaderSize(*header)
	// if err != nil {
	// 	log.Printf("Error calculating error size: %s", err)
	// }
	// hBuf := make([]byte, hSize)

	return &payloadManager{
		src:     cfg.SrcHost,
		dst:     cfg.DestHost,
		srcPort: layers.UDPPort(cfg.SrcPort),
		dstPort: layers.UDPPort(cfg.DestPort),
		srcMAC:  cfg.Interface.HardwareAddr,
		id:      uint(payId),
		freq:    fs,
		dpMap:   dps,
		header:  header,
		size:    size,
		pBuf:    payloadBuf,
		hBuf:    hBuf,
	}
}

func (pm *payloadManager) getHeader() error {
	if pm.header == nil {
		return fmt.Errorf("no header found in payload manager - ID: (%d)", pm.id)
	}

	// initialize byte index
	idx := 0

	// Iterate through each element in the header
	for i, element := range pm.header.Elements {
		if element == nil {
			return fmt.Errorf("header element at index %d is nil", i)
		}
		// Call the Bytes() method on the element to get its byte representation
		elementBytes, s, err := element.Bytes()
		if err != nil {
			// Return an informative error if getting bytes fails for any element
			return fmt.Errorf("failed to get bytes for header element at index %d: %w", i, err)
		}

		// check if there is enough room for element to be added in
		add := int(s)
		if math.MaxUint16-idx < add {
			return fmt.Errorf("not enough room in buffer for element at index %d", i)
		}

		// put data into buffer at current idx
		copy(pm.hBuf[idx:], elementBytes)

		idx += add
	}
	return nil
}

func (pm *payloadManager) buildPayload(ctx *context.Context, db *database.RedisClient) error {
	// clear out buffer
	for i := range pm.pBuf {
		pm.pBuf[i] = 0
	}

	// get payload header
	err := pm.getHeader()
	if err != nil {
		return fmt.Errorf("error retrieving payload header for ID (%s): %s", pm.id, err)
	}

	// loop through datapoints and append data by offset and size
	for id, dp := range pm.dpMap {
		d, err := db.GetData(id)
		if err != nil {
			return fmt.Errorf("error getting data point info for ID (%s): %s", pm.id, err)
		}
		// get new value
		oldVal := d["value"]
		newVal, str := dp.update(oldVal)

		// update db with new value
		d["value"] = str
		db.UpdateData(id, d)

		// append data by offset and size
		err = dp.appendData(pm.pBuf, newVal)
		if err != nil {
			return fmt.Errorf("error building data for %d: %s", id, err)
		}
	}
	return nil
}

func manager(ctx *context.Context, cfg *config.Config, cs PayloadChans, id string, pktType string, infoChan chan<- packetInfo) {
	// Set up database interface
	db := database.NewRedisClient(
		ctx,
		cfg.DbHost+":"+cfg.DbPort,
		cfg.DbPassword,
		cfg.DbNum,
		cfg.DbReadTimeout,
		cfg.DbWriteTimeout,
	)

	// get payload info
	payloadInfo, err := db.GetPayloadInfo(id)
	if err != nil {
		log.Fatalf("Did not find payload info for ID %s: %s", id, err)
	}

	// extract frequency from payload info
	f, err := strconv.ParseFloat(payloadInfo["frequency"], 64)
	if err != nil {
		log.Fatalf("Incorrect conversion of payload (%d) frequency %s Hz", id, payloadInfo["frequency"])
	}

	// convert frequency to time.Duration
	// 1 Hz = 1 second, so 1/f = seconds
	fs := time.Duration(1 / f * float64(time.Second))

	// Create new PayloadManager
	pm := newPayloadManager(cfg, id, fs, db)

	// Start processing
	for {
		select {
		case <-pm.cs.ticker.C:
			// generate new payload
			err := pm.buildPayload(ctx, db)
			if err != nil {
				pm.info.Error = true
			}
		case pkt := <-cs.writeChan:
			// Send packet
			err := pm.sendPacket(pkt)
			if err != nil {
				log.Printf("Error sending packet: %s", err)
				continue
			}

		case pkt := <-cs.readChan:
			// Process received packet
			err := processPacket(pkt)
			if err != nil {
				log.Printf("Error processing packet: %s", err)
				continue
			}
		}
	}
}
