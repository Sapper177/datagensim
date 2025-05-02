package sim

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Sapper177/datagensim/pkg/database"
)


type PacketInfo struct {
    PacketId    int
    PacketType  string
    PacketSize  int
    Direction   bool // true = tx, false = rx
    Error       bool // true = error, false = no error
    TxTime      time.Time
    ProcessTime time.Duration
}

type PayloadChans struct {
    writeChan chan<- Packet // push to write
    readChan <-chan Packet  // pull from read
    ticker *time.Ticker
}

func Sim(ctx *context.Context, cfg *Config) {
    // Set up Bus
    bus := Bus{Name: cfg.BusName, Interface: cfg.Interface.Name, IP: cfg.SrcHost.String()}

    // Set up database interface
    db := database.NewRedisClient(
        ctx,
        cfg.DbHost + ":" + cfg.DbPort,
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
    infoChan := make(chan PacketInfo, 100)

    // initialize payload routines
    initPayloads(payloadIds, db, infoChan)

    // initialize payload monitoring
    go InitMonitoring(cfg, infoChan)
}

func initPayloads(payloadIds []string, db *database.RedisClient, infoChan chan<- PacketInfo) {

    // Spawn thread for each payload
    for i := range payloadIds {

        // get payload info
        pInfo, err := db.GetPayloadInfo(payloadIds[i])
        if err != nil {
            log.Printf("No info found for ID: %s -> %s", pInfo, err)
        }
        var freq time.Duration // hz
        fmt.Sscanf(pInfo["frequency"], "%d", freq)  // get Hz in float
        freq = 1000000 / freq * time.Microsecond // in microseconds

        ticker := time.NewTicker(freq)
        defer ticker.Stop()

        // Create channels for i/o
        cs := PayloadChans {
            writeChan: make(chan Packet, 3 * len(payloadIds)),
            readChan: make(chan Packet, len(payloadIds)),
            ticker: ticker,
        }

        // spawn go routine for each payload
        go payloadManager(cs, payloadIds[i], pInfo["packet_type"], )
    }
}
