package receiver

import (
	"bytes"
	"fmt"

	"github.com/wizcas/mudever.svc/telnet/packet"
	"github.com/wizcas/mudever.svc/telnet/protocol"
)

const (
	stateNormal = byte(0)
	stateIAC    = byte(1 << 0)
	stateCmd    = byte(1 << 1)
	stateSub    = byte(1 << 2)
)

type processor struct {
	prevByte byte
	buffer   *bytes.Buffer
	value    byte
}

func newProcessor() *processor {
	return &processor{
		buffer: bytes.NewBuffer(nil),
	}
}

func (pc *processor) addState(flag byte) {
	pc.value |= flag
}

func (pc *processor) delState(flag byte) {
	pc.value &= ^flag
}

func (pc *processor) inState(flag byte) bool {
	return (pc.value & flag) == flag
}

func (pc *processor) flushData() packet.Packet {
	p := packet.NewDataPacket(pc.buffer.Bytes())
	pc.reset()
	return p
}

func (pc *processor) reset() {
	pc.buffer.Reset()
	pc.value = 0
}

func (pc *processor) proc(data []byte, chOutput chan packet.Packet, chErr chan error) {
	for _, b := range data {
		// This byte next to IAC needs be unescaped
		if pc.inState(stateIAC) {
			pc.delState(stateIAC)
			switch {
			case b == protocol.IAC: // IAC IAC (0xFF as plain data)
				pc.buffer.WriteByte(b)
			case b == protocol.SB: // IAC SB
				chOutput <- pc.flushData()
				pc.addState(stateSub)
			case b == protocol.SE: // IAC SE
				buffer := pc.buffer.Bytes()
				chOutput <- packet.NewSubPacket(buffer[0], buffer[1:])
				pc.reset()
			case b > protocol.SE && b < protocol.SB: // IAC CMD
				chOutput <- pc.flushData()
				chOutput <- packet.NewMonoCommandPacket(b)
				pc.reset()
			case b >= protocol.WILL && b <= protocol.DONT: // IAC CMD OPTION
				chOutput <- pc.flushData()
				pc.addState(stateCmd)
				// Keep this byte as command
				if err := pc.buffer.WriteByte(b); err != nil {
					chErr <- err
				}
			default:
				chOutput <- pc.flushData()
				// Invalid data
				chErr <- fmt.Errorf("Wrong data: [IAC %v]", b)
			}
		} else {
			// Leading IAC is the escape signal, and we need to examine the next byte
			// for further information
			if b == protocol.IAC {
				pc.addState(stateIAC)
				continue
			}
			// If in CMD state, this byte is the option. A complete command packet is
			// formed together with the buffered byte
			if pc.inState(stateCmd) {
				chOutput <- packet.NewOptCommandPacket(pc.buffer.Bytes()[0], b)
				pc.reset()
				continue
			}

			// For other cases, including normal data and subnegotiation data,
			// just buffer the byte until further proceeding
			if err := pc.buffer.WriteByte(b); err != nil {
				chErr <- err
			}
		}
	}
	// If we're not ending in the middle of IAC, Command or Subnegotiation,
	// Let see it as a set of complete plain data, and output.
	if pc.inState(stateNormal) && pc.buffer.Len() > 0 {
		// log.Println(pc.buffer.Bytes())
		chOutput <- pc.flushData()
	}
}
