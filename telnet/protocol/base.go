package protocol

import "fmt"

// TelnetByte represnets a byte that is meaningful to telnet transimission
type TelnetByte interface {
	fmt.Stringer
	Eq(b byte) bool
}

// CmdByte represents a command byte in telnet transimission
type CmdByte byte

func (cmd CmdByte) String() string {
	name, ok := cmdNames[cmd]
	if !ok {
		return fmt.Sprintf("%v", byte(cmd))
	}
	return name
}

// Eq is a short syntax to check the equality to a given byte
func (cmd CmdByte) Eq(b byte) bool {
	return byte(cmd) == b
}

// OptByte represents an option byte in telnet transimission
type OptByte byte

func (opt OptByte) String() string {
	name, ok := optionNames[opt]
	if !ok {
		return fmt.Sprintf("%v", byte(opt))
	}
	return name
}

// Eq is a short syntax to check the equality to a given byte
func (opt OptByte) Eq(b byte) bool {
	return byte(opt) == b
}
