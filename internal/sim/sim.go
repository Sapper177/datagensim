package sim

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Sapper177/datagensim/pkg/config"
	"github.com/Sapper177/datagensim/pkg/database"
)

type PayloadChans struct {
	writeChan chan Packet // push to write
	readChan  chan Packet // pull from read
	ticker    *time.Ticker
}

func Sim(ctx *context.Context, cfg *config.Config) {
	// Set up Bus
	// bus := Bus{Name: cfg.BusName, Interface: cfg.Interface.Name, IP: cfg.SrcHost.String()}

	// Set up database interface
	db := database.NewRedisClient(
		ctx,
		cfg.DbHost+":"+cfg.DbPort,
		cfg.DbPassword,
		cfg.DbNum,
		cfg.DbReadTimeout,
		cfg.DbWriteTimeout,
	)

	// Get payload configs from database
	payloadIds, err := db.GetPayloads(cfg.BusName)
	if err != nil {
		log.Fatalf("Did not find payload IDs for Bus %s: %s", cfg.BusName, err)
	}

	// Create channel that will be used contain sent packet data
	infoChan := make(chan packetInfo, 100)

	// initialize payload routines
	initPayloads(ctx, cfg, payloadIds, db, infoChan)

	// initialize payload monitoring
	go initMonitoring(cfg, infoChan)

	// Run Simulation
	// go sim(payloadManagers, infoChan)
}

func initPayloads(ctx *context.Context, cfg *config.Config, payloadIds []string, db *database.RedisClient, infoChan chan<- packetInfo) {

	// Spawn thread for each payload
	for i := range payloadIds {

		// get payload info
		pInfo, err := db.GetPayloadInfo(payloadIds[i])
		if err != nil {
			log.Printf("No info found for ID: %s -> %s", pInfo, err)
		}
		var freq time.Duration                     // hz
		fmt.Sscanf(pInfo["frequency"], "%d", freq) // get Hz in float
		freq = 1000000 / freq * time.Microsecond   // in microseconds

		ticker := time.NewTicker(freq)
		defer ticker.Stop()

		// Create channels for i/o
		cs := PayloadChans{
			writeChan: make(chan Packet, 3*len(payloadIds)),
			readChan:  make(chan Packet, len(payloadIds)),
			ticker:    ticker,
		}

		// spawn go routine for each payload
		go manager(ctx, cfg, cs, payloadIds[i], pInfo["packet_type"], infoChan)
	}
}

// func sim(payloadManagers []*payloadManager, infoChan chan<- packetInfo) {
// 	for _, pm := range payloadManagers {

// 		// check if pm is ready to send from frequency
// 		if pm != nil {

// 			if time.Since(pm.lastProc) >= pm.freq {
// 				// generate packet
// 				pkt, err := pm.generatePacket(infoChan)
// 				intId, err1 := strconv.Atoi(pm.id)
// 				if err1 != nil {
// 					log.Printf("Error converting payload ID to int: %s", err)
// 				}

// 				if err != nil {
// 					log.Printf("Error generating packet: %s", err)

// 				} else {
// 					pm.cs.writeChan <- pkt   // send packet to write channel
// 					pm.lastProc = time.Now() // update last processed time
// 					// send packet info to info channel
// 					infoChan <- packetInfo{
// 						PacketId:    intId,
// 						PacketType:  pm.pktType,
// 						PacketSize:  pkt.Size(),
// 						Direction:   true,
// 						Error:       false,
// 						TxTime:      time.Now(),
// 						ProcessTime: procTime,
// 					}
// 				}

// 			} // wait until process time is up

// 		} else {
// 			log.Printf("Payload manager is nil for ID: %s", pm.id)
// 		}
// 	}
// }
