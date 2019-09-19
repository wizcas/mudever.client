package packet

// Packet is the minimum processing unit for telnet transimission
type Packet interface {
	// Serialize the packet into the form ready for transimission.
	// E.g., it escape the command byte with IAC
	Serialize() ([]byte, error)
}
