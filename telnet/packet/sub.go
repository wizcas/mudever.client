package packet

import (
	"bytes"
	"fmt"

	"github.com/wizcas/mudever.svc/telnet/protocol"
)

// SubPacket represents a telnet subnegotiation
type SubPacket struct {
	Option    byte
	Parameter []byte
}

// NewSubPacket create subnegotiation message with given option and optional args.
// 'args' will be concat as the PARAMETER between OPTION & SE
func NewSubPacket(option byte, args ...[]byte) *SubPacket {
	parameter := []byte{}
	for _, arg := range args {
		parameter = append(parameter, arg...)
	}
	return &SubPacket{option, parameter}
}

// Serialize the packet with telnet protocol, i.e.:
// IAC SB <OPTION> <PARAMETER> IAC SE
func (p *SubPacket) Serialize() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{protocol.IAC, protocol.SB, p.Option})
	if n, err := buf.Write(p.Parameter); err != nil {
		return nil, err
	} else if n != len(p.Parameter) {
		return nil, err
	}
	if err := buf.WriteByte(protocol.IAC); err != nil {
		return nil, err
	}
	if err := buf.WriteByte(protocol.SE); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *SubPacket) String() string {
	return fmt.Sprintf("[SUB] %v | %v", p.Option, p.Parameter)
}
