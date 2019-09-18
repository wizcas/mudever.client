package cmds

// Standard Commands
const (
	// SE = Subnegotiation ends
	SE = byte(240)
	// NOP = No operation
	NOP = byte(241)
	// DM = Data mark (pos of Synch event)
	DM = byte(242)
	// BRK = Break
	BRK = byte(243)
	// IP = Interrupt the process to which the NVT is connected
	IP = byte(244)
	// AO = Abort output
	AO = byte(245)
	// AYT = Are you there
	AYT = byte(246)
	// EC = Erase character
	EC = byte(247)
	// EL = Erase line excluding the previous CRLF
	EL = byte(248)
	// GA = Go ahead
	GA = byte(249)
	// SB = Subnegotiation begins
	SB = byte(250)
	// WILL = Sender wants to do something
	WILL = byte(251)
	// WONT = Sender doesn't to do something
	WONT = byte(252)
	// DO = Sender wants THE OTHER to do something
	DO = byte(253)
	// DONT = Sender doesn't THE OTHER to do something
	DONT = byte(254)
	// IAC = Interpret as command. The escape character in telnet protocol.
	IAC = byte(255)
)

// Standard Options
const (
	Echo = byte(1)
	SuppressGoAhead = byte(3)
	Status = byte(5)
	TimingMark = byte(6)
	TerminalType = byte(24)
	WindowSize = byte(31)
	TerminalSpeed = byte(32)
	RemoteFlowControl = byte(33)
	LineMode = byte(34)
	EnvironmentVariables = byte(36)
)