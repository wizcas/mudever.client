package nego

import (
	"errors"

	"github.com/wizcas/mudever.svc/packet"
	"github.com/wizcas/mudever.svc/telbyte"
)

// Handler is registered to negotiator for telnet command interpretion
type Handler interface {
}

// ControlHandler takes care of control commands (e.g. GA, IP, etc.) in its Handle function.
// Note that the option commands (i.e. WILL, WONT, DO, DONT) should not be handled by ControlHandler,
// but rather to be processed within an OptionHandler.
type ControlHandler interface {
	Handler
	Command() telbyte.Command
	Handle() error
}

// HandlerCommittee takes any result or error that committed by a handler.
type HandlerCommittee interface {
	Commit(p packet.Packet) error
	GotError(err error)
}

// Errors caused by handlers
var (
	ErrLackData = errors.New("LACK OF DATA")
)
