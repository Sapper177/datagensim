package sim

type Packet struct {
    Protocol string // "udp" or "tcp"
    SrcPort  int
    DstPort  int
    Payload  []byte
}
