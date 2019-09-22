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

type OptionContext struct {
	context.Context
	Cancel  context.CancelFunc
	handler OptionHandler
	sender  common.PacketSender
	onError common.OnError
	Logger  *zap.SugaredLogger
}

func newOptionContext(parentCtx context.Context, handler OptionHandler, ng *Negotiator) *OptionContext {
	ctx, cancel := context.WithCancel(parentCtx)
	return &OptionContext{
		Context: ctx,
		Cancel:  cancel,
		handler: handler,
		sender:  ng.sender,
		onError: ng.GotError,
		Logger:  logger().Named(handler.Option().String()),
	}
}

func (c *OptionContext) GotError(err error) {
	c.onError(err)
}
func (c *OptionContext) SendCmd(cmd telbyte.Command) {
	p := packet.NewOptionCommandPacket(cmd, c.handler.Option())
	if err := c.sender.Send(p); err != nil {
		c.GotError(err)
	}
}
func (c *OptionContext) SendSub(parameters ...[]byte) {
	p := packet.NewSubPacket(c.handler.Option(), parameters...)
	if err := c.sender.Send(p); err != nil {
		c.GotError(err)
	}
}
