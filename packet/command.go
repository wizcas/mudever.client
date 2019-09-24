package packet

import (
	"fmt"

	"github.com/wizcas/mudever.svc/telbyte"
)

// CommandPacket represents a telnet command
type CommandPacket struct {
	Command telbyte.Command
	Option  telbyte.Option
}

// NewControlCommandPacket creates a command of a telnet control function
func NewControlCommandPacket(command telbyte.Command) *CommandPacket {
	return &CommandPacket{command, telbyte.NoOption}
}

// NewOptionCommandPacket creates a command of a negotiable option
func NewOptionCommandPacket(command telbyte.Command, option telbyte.Option) *CommandPacket {
	return &CommandPacket{command, option}
}

// Serialize the packet with IAC-escape,
// which makes it contain 2 bytes for mono command and 3 for option command.
func (p *CommandPacket) Serialize() ([]byte, error) {
	bytes := [3]byte{byte(telbyte.IAC), byte(p.Command), byte(p.Option)}
	if !p.IsOption() {
		return bytes[:2], nil
	}
	return bytes[:], nil
}

// IsOption indicates whether this command packet refers
// an option command or a control command
func (p *CommandPacket) IsOption() bool {
	return p.Option != telbyte.NoOption
}

func (p *CommandPacket) String() string {
	var str string
	if !p.IsOption() {
		str = fmt.Sprintf("%s", p.Command)
	} else {
		str = fmt.Sprintf("%s > %s", p.Command, p.Option)
	}
	return fmt.Sprintf("[CMD] %s", str)
}
