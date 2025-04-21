package sim

import (
	"context"
	"gosim/pkg/database"
	"log"
	"time"

	"pkg/database"
)


type PacketInfo struct {
    PacketId    int
    PacketType  string
    PacketSize  int
    TxTime      time.Time
    ProcessTime time.Time
}

func Sim(ctx *context.Context, cfg *Config) {
    // Set up Bus
    bus := Bus{Name: cfg.BusName, Interface: cfg.Interface.Name, IP: cfg.SrcHost.String()}

    // Set up database interface
    db := database.NewRedisClient(
        cfg.DbHost + ":" + cfg.DbPort,
        cfg.DbPassword,
        cfg.DbNum,
        cfg.DbReadTimeout,
        cfg.DbWriteTimeout,
    )

    // Get payload configs from database
    payloadConfigs, err := 

    // Create channel that will be used contain sent packet data
    infoChan := make(chan PacketInfo, 100)

    // Create channel for sending packets
    packetChan := make(chan Packet, 50)

    // Spawn thread for each payload




    // packet := types.Packet{Protocol: "udp", SrcPort: 12345, DstPort: 54321}
    // dataItems := []types.DataItem{
    //     {Name: "Temperature", Type: types.Ramp, Value: 20.0, Min: 20.0, Max: 30.0, Step: 0.5},
    //     {Name: "Pressure", Type: types.Sinusoidal, Value: 1.0, Min: 0.5, Max: 1.5, Frequency: 0.1, Phase: 0},
    //     {Name: "Status", Type: types.Static, Value: 1},
    // }

    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        for i := range dataItems {
            dataItems[i].Update()
        }
        payload := constructPayload(dataItems)
        packet.Payload = payload
        sendPacket(bus, packet)
    }
}
