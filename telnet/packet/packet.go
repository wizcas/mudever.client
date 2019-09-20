package packet

import "fmt"

// Packet is the minimum processing unit for telnet transimission
type Packet interface {
	fmt.Stringer
	// Serialize the packet into the form ready for transimission.
	// E.g., it escape the command byte with IAC
	Serialize() ([]byte, error)
}
