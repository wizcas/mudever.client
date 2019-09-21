package negotiator

import (
	"errors"

	"github.com/wizcas/mudever.svc/telnet/protocol"
)

// Handler is registered to negotiator for telnet command interpretion
type Handler interface{}

// ControlHandler takes care of control commands (e.g. GA, IP, etc.) in its Handle function.
// Note that the option commands (i.e. WILL, WONT, DO, DONT) should not be handled by ControlHandler,
// but rather to be processed within an OptionHandler.
type ControlHandler interface {
	Handler
	Command() protocol.CmdByte
	Handle() error
}

// OptionHandler takes care of option commands and subnegotiations of a certain Telnet Option
type OptionHandler interface {
	Handler
	Option() protocol.OptByte
	Handshake(inCmd protocol.CmdByte) (protocol.CmdByte, error)
	Subnegotiate(inParameter []byte) ([]byte, error)
}

// Errors caused by handlers
var (
	ErrIgnore   = errors.New("IGNORE")
	ErrLackData = errors.New("LACK OF DATA")
)
