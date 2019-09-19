package packet

import (
	"fmt"

	"github.com/wizcas/mudever.svc/telnet/protocol"
)

// CommandPacket represents a telnet command
type CommandPacket struct {
	Command byte
	Option  byte
}

// NewMonoCommandPacket creates a no-option command.
func NewMonoCommandPacket(command byte) *CommandPacket {
	return &CommandPacket{command, protocol.IAC}
}

// NewOptCommandPacket creates a command with an option.
func NewOptCommandPacket(command byte, option byte) *CommandPacket {
	return &CommandPacket{command, option}
}

// Serialize the packet with IAC-escape,
// which makes it contain 2 bytes for mono command and 3 for option command.
func (p *CommandPacket) Serialize() ([]byte, error) {
	bytes := [3]byte{protocol.IAC, p.Command, p.Option}
	if p.isMono() {
		return bytes[:2], nil
	}
	return bytes[:], nil
}

func (p *CommandPacket) isMono() bool {
	return p.Option == protocol.IAC
}

func (p *CommandPacket) String() string {
	var str string
	if p.isMono() {
		str = fmt.Sprintf("%v", p.Command)
	} else {
		str = fmt.Sprintf("%v | %v", p.Command, p.Option)
	}
	return fmt.Sprintf("[CMD] %s", str)
}
