package receiver

import (
	"bytes"
	"fmt"

	"github.com/wizcas/mudever.svc/telnet/nvt/common"
	"github.com/wizcas/mudever.svc/telnet/packet"
	"github.com/wizcas/mudever.svc/telnet/telbyte"
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

func (pc *processor) flush() []byte {
	buffered := pc.buffer.Bytes()
	data := make([]byte, len(buffered), len(buffered))
	copy(data, buffered)
	return data
}

func (pc *processor) submitAsData(chOutput chan packet.Packet) {
	data := pc.flush()
	if len(data) > 0 {
		p := packet.NewDataPacket(data)
		chOutput <- p
	}
	pc.reset()
}

func (pc *processor) reset() {
	pc.buffer.Reset()
	pc.value = 0
}

func (pc *processor) proc(data []byte, chOutput chan packet.Packet, onError common.OnError) {
	for _, b := range data {
		// log.Printf("%v\n", b)
		// This byte next to IAC needs be unescaped
		if pc.inState(stateIAC) {
			pc.delState(stateIAC)
			switch {
			case telbyte.IAC.Eq(b): // IAC IAC (0xFF as plain data)
				pc.buffer.WriteByte(b)
			case telbyte.SB.Eq(b): // IAC SB
				pc.submitAsData(chOutput)
				pc.addState(stateSub)
			case telbyte.SE.Eq(b): // IAC SE
				buffer := pc.flush()
				chOutput <- packet.NewSubPacket(telbyte.Option(buffer[0]), buffer[1:])
				pc.reset()
			case b > byte(telbyte.SE) && b < byte(telbyte.SB): // IAC CMD
				pc.submitAsData(chOutput)
				chOutput <- packet.NewControlCommandPacket(telbyte.Command(b))
				pc.reset()
			case b >= byte(telbyte.WILL) && b <= byte(telbyte.DONT): // IAC CMD OPTION
				pc.submitAsData(chOutput)
				pc.addState(stateCmd)
				// Keep this byte as command
				if err := pc.buffer.WriteByte(b); err != nil {
					onError(err)
				}
			default:
				pc.submitAsData(chOutput)
				// Invalid data
				onError(fmt.Errorf("Wrong data: [IAC %v]", b))
			}
		} else {
			// Leading IAC is the escape signal, and we need to examine the next byte
			// for further information
			if telbyte.IAC.Eq(b) {
				// log.Println("iac")
				pc.addState(stateIAC)
				continue
			}
			// If in CMD state, this byte is the option. A complete command packet is
			// formed together with the buffered byte
			if pc.inState(stateCmd) {
				// log.Println("opt cmd ends")
				buffer := pc.flush()
				chOutput <- packet.NewOptionCommandPacket(telbyte.Command(buffer[0]), telbyte.Option(b))
				pc.reset()
				continue
			}

			// For other cases, including normal data and subnegotiation data,
			// just buffer the byte until further proceeding
			if err := pc.buffer.WriteByte(b); err != nil {
				onError(err)
			}
		}
	}
	// If we're not ending in the middle of IAC, Command or Subnegotiation,
	// Let see it as a set of complete plain data, and output.
	if pc.inState(stateNormal) && pc.buffer.Len() > 0 {
		// log.Println(pc.buffer.Bytes())
		pc.submitAsData(chOutput)
	}
}
