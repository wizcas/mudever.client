package packet

// Kind of the telnet packet
type Kind byte

const (
	// KindData indicates that the packet should be interpreted as USASCII characters
	KindData = Kind(iota)
	// KindCommand indicates that the packet is a telnet option command
	KindCommand
	// KindSubnegotiation indicates that the packet is a telnet subnegotiation
	KindSubnegotiation
)

// Packet is the minimum processing unit for telnet transimission
type Packet interface {
	// GetKind is used for determine which kind of packet this is
	GetKind() Kind
	// Serialize the packet into the form ready for transimission.
	// E.g., it escape the command byte with IAC
	Serialize() ([]byte, error)
}
