package nego

import (
	"context"

	"github.com/wizcas/mudever.svc/nvt/common"
	"github.com/wizcas/mudever.svc/packet"
	"github.com/wizcas/mudever.svc/telbyte"
	"go.uber.org/zap"
)

// OptionHandler takes care of option commands and subnegotiations of a certain Telnet Option
type OptionHandler interface {
	Handler
	Option() telbyte.Option
	Handshake(ctx *OptionContext, inCmd telbyte.Command)
	Subnegotiate(ctx *OptionContext, inParameter []byte)
}

// OptionContext passes necessary callbacks and objects, such as logger, context, etc., to
// option handlers. Handler should call SendCmd() or SendSub() to write data back to terminal,
// and call GotError() to notify terminal on error occurance. The context is cancelled when
// Cancel() is called manually or when its parent context is cancelled.
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

// NewOptionContext returns a new option context
func NewOptionContext(parentCtx context.Context, handler OptionHandler,
	sender common.PacketSender,
	onError common.OnError,
) *OptionContext {
	ctx, cancel := context.WithCancel(parentCtx)
	return &OptionContext{
		Context: ctx,
		Cancel:  cancel,
		handler: handler,
		sender:  sender,
		onError: onError,
		Logger:  common.Logger().With("handler", handler.Option()),
	}
}

// GotError should be called to report errors in handler to terminal.
func (c *OptionContext) GotError(err error) {
	c.onError(err)
}

// SendCmd takes a command byte, packed into an OptionCommandPacket, and
// sent to server as a telnet option command.
func (c *OptionContext) SendCmd(cmd telbyte.Command) {
	p := packet.NewOptionCommandPacket(cmd, c.handler.Option())
	if err := c.sender.Send(p); err != nil {
		c.GotError(err)
	}
}

// SendSub takes subnegotiation parameter(s), packed into a SubPacket, and
// send to server as a telnet subnegotiation.
func (c *OptionContext) SendSub(parameters ...[]byte) {
	p := packet.NewSubPacket(c.handler.Option(), parameters...)
	if err := c.sender.Send(p); err != nil {
		c.GotError(err)
	}
}
