package nego

import (
	"context"
	"errors"

	"github.com/wizcas/mudever.svc/telnet/nvt/common"
	"github.com/wizcas/mudever.svc/telnet/packet"
	"github.com/wizcas/mudever.svc/telnet/telbyte"
	"go.uber.org/zap"
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

// OptionHandler takes care of option commands and subnegotiations of a certain Telnet Option
type OptionHandler interface {
	Handler
	Option() telbyte.Option
	Handshake(ctx *OptionContext, inCmd telbyte.Command)
	Subnegotiate(ctx *OptionContext, inParameter []byte)
}

// Errors caused by handlers
var (
	ErrLackData = errors.New("LACK OF DATA")
)

// OptionContext passes necessary callbacks and objects, such as logger, context, etc., to
// option handlers. Handler should call SendCmd() or SendSub() to write data back to terminal,
// and call GotError() to notify terminal on error occurance. The context is cancelled when
// Cancel() is called manually or its creator's (i.e. negotiator's) context is cancelled.
type OptionContext struct {
	context.Context
	handler OptionHandler
	sender  common.PacketSender
	onError common.OnError
	// Cancel the handler's context when needed.
	Cancel context.CancelFunc
	// Logger provides a main.nvt logger with the option's name as the value of 'handler' field.
	Logger *zap.SugaredLogger
}

func newOptionContext(parentCtx context.Context, handler OptionHandler, ng *Negotiator) *OptionContext {
	ctx, cancel := context.WithCancel(parentCtx)
	return &OptionContext{
		Context: ctx,
		Cancel:  cancel,
		handler: handler,
		sender:  ng.sender,
		onError: ng.GotError,
		Logger:  common.Logger().With("handler", handler.Option()),
	}
}

// GotError should be called to report errors in handler to terminal.
func (c *OptionContext) GotError(err error) {
	c.onError(err)
}

// SendCmd sends a command byte to terminal, which will be encoded and
// sent to server as a telnet option command.
func (c *OptionContext) SendCmd(cmd telbyte.Command) {
	p := packet.NewOptionCommandPacket(cmd, c.handler.Option())
	if err := c.sender.Send(p); err != nil {
		c.GotError(err)
	}
}

// SendSub sends subnegotiation parameter(s) to terminal, which will be
// encoded and send to server as a telnet subnegotiation.
func (c *OptionContext) SendSub(parameters ...[]byte) {
	p := packet.NewSubPacket(c.handler.Option(), parameters...)
	if err := c.sender.Send(p); err != nil {
		c.GotError(err)
	}
}
