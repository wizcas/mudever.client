package packet

import (
	"github.com/wizcas/mudever.svc/telnet/protocol"
)

// DataPacket represents a bunch of telnet plain data
type DataPacket struct {
	data []byte
}

// NewDataPacket creates a data packet with given data
func NewDataPacket(data []byte) *DataPacket {
	return &DataPacket{data}
}

// GetKind returns KindData
func (p *DataPacket) GetKind() Kind {
	return KindData
}

// Serialize the packet with all 0xFF character escaped by IAC
func (p *DataPacket) Serialize() []byte {
	result := p.data[:]
	for i := 0; i < len(result); i++ {
		b := result[i]
		if b == protocol.IAC {
			result = insert(result, protocol.IAC, i)
			i++
		}
	}
	return result
}

func insert(dst []byte, b byte, pos int) []byte {
	dst = append(dst, 0)
	copy(dst[pos+1:], dst[pos:])
	dst[pos] = b
	return dst
}

func insertRef(dst *[]byte, b byte, pos int) {
	slice := *dst
	slice = append(slice, 0)
	copy(slice[pos+1:], slice[pos:])
	slice[pos] = b
	*dst = slice
}
