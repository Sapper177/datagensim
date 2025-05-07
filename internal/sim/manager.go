package sim

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/Sapper177/datagensim/internal/engine"
	"github.com/Sapper177/datagensim/internal/pktgen"
	"github.com/Sapper177/datagensim/pkg/database"
	"github.com/google/gopacket/layers"
)
type PayloadManager struct {
	pktType string // UDP, TCP, etc
	src net.IP
	dst net.IP
	srcPort layers.UDPPort
	dstPort layers.UDPPort
	srcMAC net.HardwareAddr
	id string		// payload id
	freq time.Duration
	lastProc time.Time
	cs *PayloadChans
}
func NewPayloadManager(cfg *Config, id string, fs time.Duration) *PayloadManager{
	return &PayloadManager{
		src: cfg.SrcHost,
		dst: cfg.DestHost,
		srcPort: layers.UDPPort(cfg.SrcPort),
		dstPort: layers.UDPPort(cfg.DestPort),
		srcMAC: cfg.Interface.HardwareAddr,
		id: id,
		freq: fs,
	}
}

func payloadManager(ctx *context.Context, cfg *Config, cs PayloadChans, id string, pktType string, pmlist []*PayloadManager) {
	// Set up database interface
	db := database.NewRedisClient(
		ctx,
		cfg.DbHost + ":" + cfg.DbPort,
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
	pm := NewPayloadManager(cfg, id, fs)
	pmlist = append(pmlist, pm) // add to pm list

	// Start thread for packet generation
	for {
		select {
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