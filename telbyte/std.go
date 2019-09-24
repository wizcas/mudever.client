package telbyte

// Standard Commands
const (
	// SE = Subnegotiation ends
	SE = Command(240)
	// NOP = No operation
	NOP = Command(241)
	// DM = Data mark (pos of Synch event)
	DM = Command(242)
	// BRK = Break
	BRK = Command(243)
	// IP = Interrupt the process to which the NVT is connected
	IP = Command(244)
	// AO = Abort output
	AO = Command(245)
	// AYT = Are you there
	AYT = Command(246)
	// EC = Erase character
	EC = Command(247)
	// EL = Erase line excluding the previous CRLF
	EL = Command(248)
	// GA = Go ahead
	GA = Command(249)
	// SB = Subnegotiation begins
	SB = Command(250)
	// WILL = Sender wants to do something
	WILL = Command(251)
	// WONT = Sender doesn't to do something
	WONT = Command(252)
	// DO = Sender wants THE OTHER to do something
	DO = Command(253)
	// DONT = Sender doesn't THE OTHER to do something
	DONT = Command(254)
	// IAC = Interpret as command. The escape character in telnet telbyte.
	IAC = Command(255)
)

// Standard Options
const (
	// Echo
	ECHO = Option(1)
	// Suppress Go Ahead
	SUPPRESSGA = Option(3)
	// Status
	STATUS = Option(5)
	// Timing Mark
	TM = Option(6)
	// Terminal Type
	TTYPE = Option(24)
	// Negotiate About Window Size
	NAWS = Option(31)
	// Terminal Speed
	TSPEED = Option(32)
	// Toggle Flow Control / Remote Flow Control
	TOGFC = Option(33)
	// Line Mode
	LINEMODE = Option(34)
	// Environment Variables
	ENVIRON = Option(36)
	// An invalid value for option byte
	NoOption = Option(255)
)
