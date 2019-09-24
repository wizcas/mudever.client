package packet

import (
	"fmt"
)

// DataPacket represents a bunch of telnet plain data
type DataPacket struct {
	// Data to be output
	Data []byte
}

// NewDataPacket creates a data packet with given data
func NewDataPacket(data []byte) *DataPacket {
	return &DataPacket{data}
}

// Serialize the packet with all 0xFF character escaped by IAC
func (p *DataPacket) Serialize() ([]byte, error) {
	return escapeData(p.Data), nil
}

func (p *DataPacket) String() string {
	return fmt.Sprintf("[TXT] (%d bytes)", len(p.Data))
}
