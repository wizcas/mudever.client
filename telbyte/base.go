package telbyte

import "fmt"

// TelByte represnets a byte that is meaningful to telnet transimission
type TelByte interface {
	fmt.Stringer
	Eq(b byte) bool
}

// Command represents a command byte in telnet transimission
type Command byte

func (cmd Command) String() string {
	name, ok := cmdNames[cmd]
	if !ok {
		return fmt.Sprintf("%v", byte(cmd))
	}
	return name
}

// Eq is a short syntax to check the equality to a given byte
func (cmd Command) Eq(b byte) bool {
	return byte(cmd) == b
}

// Option represents an option byte in telnet transimission
type Option byte

func (opt Option) String() string {
	name, ok := optionNames[opt]
	if !ok {
		return fmt.Sprintf("%v", byte(opt))
	}
	return name
}

// Eq is a short syntax to check the equality to a given byte
func (opt Option) Eq(b byte) bool {
	return byte(opt) == b
}
