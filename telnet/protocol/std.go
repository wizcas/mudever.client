package protocol

// Standard Commands
const (
	// SE = Subnegotiation ends
	SE = CmdByte(240)
	// NOP = No operation
	NOP = CmdByte(241)
	// DM = Data mark (pos of Synch event)
	DM = CmdByte(242)
	// BRK = Break
	BRK = CmdByte(243)
	// IP = Interrupt the process to which the NVT is connected
	IP = CmdByte(244)
	// AO = Abort output
	AO = CmdByte(245)
	// AYT = Are you there
	AYT = CmdByte(246)
	// EC = Erase character
	EC = CmdByte(247)
	// EL = Erase line excluding the previous CRLF
	EL = CmdByte(248)
	// GA = Go ahead
	GA = CmdByte(249)
	// SB = Subnegotiation begins
	SB = CmdByte(250)
	// WILL = Sender wants to do something
	WILL = CmdByte(251)
	// WONT = Sender doesn't to do something
	WONT = CmdByte(252)
	// DO = Sender wants THE OTHER to do something
	DO = CmdByte(253)
	// DONT = Sender doesn't THE OTHER to do something
	DONT = CmdByte(254)
	// IAC = Interpret as command. The escape character in telnet protocol.
	IAC = CmdByte(255)
)

// Standard Options
const (
	Echo                 = OptByte(1)
	SuppressGoAhead      = OptByte(3)
	Status               = OptByte(5)
	TimingMark           = OptByte(6)
	TerminalType         = OptByte(24)
	NAWS                 = OptByte(31)
	TerminalSpeed        = OptByte(32)
	RemoteFlowControl    = OptByte(33)
	LineMode             = OptByte(34)
	EnvironmentVariables = OptByte(36)

	NoOption = OptByte(255)
)
