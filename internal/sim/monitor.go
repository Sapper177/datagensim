package sim

import (
	"time"
    "sync"
)

type PayloadInfo struct {
    txCount int
    txMbs float64
    rxCount int
    rxMbs float64
    errCount int
    errMbs float64
    avgProcessingTime time.Duration
}
func (p *PayloadInfo) updateTxMbs(numBytes int) {
    p.rxMbs += float64(numBytes / 1000)
}
func (p *PayloadInfo) updateRxMbs(numBytes int) {
    p.txMbs += float64(numBytes / 1000)
}
func (p *PayloadInfo) updateErrMbs(numBytes int) {
    p.errMbs += float64(numBytes / 1000)
}
func (p *PayloadInfo) updateAvgProcessingTime(processingTime time.Duration) {
    p.avgProcessingTime = (processingTime + p.avgProcessingTime) / 2
}

type PayloadMonitor struct {
    mu sync.RWMutex
    payloadMap map[string]PayloadInfo
    numPayloads int
    numTx int
    numRx int
    numErr int
    numMbsTx float64
    numMbsRx float64
    numErrMbs float64
    totalProcessingTime time.Duration
    processedOperations int
}
// NewPayloadMonitor creates a new instance of the monitor.
func NewPayloadMonitor() *PayloadMonitor {
	return &PayloadMonitor{
		payloadMap: make(map[string]PayloadInfo),
	}
}
func (p *PayloadMonitor) addInfo(info PacketInfo) {
    p.mu.Lock()
	defer p.mu.Unlock()
    // check if payload exists in map
    id := string(info.PacketId)
	payInfo, exists := p.payloadMap[id]
	if !exists {
        // if not, create a new PayloadInfo
        newInfo := PayloadInfo{
            txCount: 0,
            txMbs: 0,
            rxCount: 0,
            rxMbs: 0,
            errCount: 0,
            errMbs: 0,
            avgProcessingTime: 0,
        }
		p.payloadMap[id] = newInfo
        payInfo = newInfo
		p.numPayloads++
	}
    p.totalProcessingTime += info.ProcessTime
	p.processedOperations++
    if info.Error { // if packet had an error
        payInfo.updateErrMbs(info.PacketSize)
        p.numErr++
        p.numErrMbs += payInfo.errMbs
    } else {
        payInfo.updateAvgProcessingTime(info.ProcessTime)   // update average processing time in packetInfo
        if info.Direction { // if packet is tx
            payInfo.updateTxMbs(info.PacketSize)
            p.numTx++
            p.numMbsTx += payInfo.txMbs
        } else {    // if packet is rx
            payInfo.updateRxMbs(info.PacketSize)
            p.numRx++
            p.numMbsRx += payInfo.rxMbs
        }
    }
}

func procPayloadMon(cfg *Config, payloadMon *PayloadMonitor, infoChan <-chan PacketInfo) {
    // process PacketInfos
    for {
        select {
        case info := <-infoChan:
            payloadMon.addInfo(info)
        case <-time.After(cfg.MonitorInterval):
            // do nothing
        }
    }
}