package nego

import (
	"errors"

	"github.com/wizcas/mudever.svc/telnet/telbyte"
)

// Handler is registered to negotiator for telnet command interpretion
type Handler interface {
	// Initialize should store the passed-in negotiator and use it for further message delivery,
	// via its chHandled* channels
	Initialize(nego *Negotiator)
}

// ControlHandler takes care of control commands (e.g. GA, IP, etc.) in its Handle function.
// Note that the option commands (i.e. WILL, WONT, DO, DONT) should not be handled by ControlHandler,
// but rather to be processed within an OptionHandler.
type ControlHandler interface {
	Handler
	Command() telbyte.Command
	Handle() error
}

// OptionHandler takes care of option commands and subnegotiations of a certain Telnet Option
type OptionHandler interface {
	Handler
	Option() telbyte.Option
	Handshake(inCmd telbyte.Command)
	Subnegotiate(inParameter []byte)
}

// Errors caused by handlers
var (
	ErrIgnore   = errors.New("IGNORE")
	ErrLackData = errors.New("LACK OF DATA")
)

type HandledCommand struct {
	Handler OptionHandler
	Command telbyte.Command
}

func NewHandledCmd(handler OptionHandler, cmd telbyte.Command) HandledCommand {
	return HandledCommand{
		Handler: handler,
		Command: cmd,
	}
}

type HandledSub struct {
	Handler   OptionHandler
	Parameter [][]byte
}

func NewHandledSub(handler OptionHandler, parameters ...[]byte) HandledSub {
	return HandledSub{
		Handler:   handler,
		Parameter: parameters,
	}
}

type OptionHandlerBase struct {
	ChOutCmd chan HandledCommand
	ChOutSub chan HandledSub
	ChErr    chan error
}

func NewOptionHandlerBase() *OptionHandlerBase {
	return &OptionHandlerBase{}
}
func (h *OptionHandlerBase) Initialize(nego *Negotiator) {
	h.ChOutCmd = nego.ChHandledCmd
	h.ChOutSub = nego.ChHandledSub
	h.ChErr = nego.ChErr
}
