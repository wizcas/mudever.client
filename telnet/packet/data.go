package packet

import (
	"fmt"

	"github.com/wizcas/mudever.svc/telnet/protocol"
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
	result := p.Data[:]
	for i := 0; i < len(result); i++ {
		b := result[i]
		if b == protocol.IAC {
			result = insert(result, protocol.IAC, i)
			i++
		}
	}
	return result, nil
}

func (p *DataPacket) String() string {
	return fmt.Sprintf("[TXT] (%d bytes)", len(p.Data))
}

func insert(dst []byte, b byte, pos int) []byte {
	dst = append(dst, 0)
	copy(dst[pos+1:], dst[pos:])
	dst[pos] = b
	return dst
}
