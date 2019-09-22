package nego

import (
	"context"

	"github.com/wizcas/mudever.svc/telnet/nvt/log"

	"github.com/wizcas/mudever.svc/telnet/nvt/common"
	"github.com/wizcas/mudever.svc/telnet/packet"
	"github.com/wizcas/mudever.svc/telnet/telbyte"
	"go.uber.org/zap"
)

// Negotiator takes care of telnet negotiations
type Negotiator struct {
	*common.SubProc
	controlHandlers map[telbyte.Command]ControlHandler
	optionHandlers  map[telbyte.Option]OptionHandler

	chInput chan packet.Packet
	sender  common.PacketSender
}

func logger() *zap.SugaredLogger {
	return log.Logger().Named("nego")
}

// New negotiator takes a PacketSender for its handlers to feed replies.
// It is created with an empty knowledge base, i.e. you'll have to feed in
// Handlers with Know() method.
func New(sender common.PacketSender) *Negotiator {
	return &Negotiator{
		SubProc:         common.NewSubProc(),
		optionHandlers:  make(map[telbyte.Option]OptionHandler),
		controlHandlers: make(map[telbyte.Command]ControlHandler),
		chInput:         make(chan packet.Packet),
		sender:          sender,
	}
}

// Know an handler makes the negotiator being able to deal with the corresponding incoming option
func (nego *Negotiator) Know(handler Handler) {
	switch h := handler.(type) {
	case ControlHandler:
		logger().Debugf("Control Handler: %s", h.Command())
		nego.controlHandlers[h.Command()] = h
	case OptionHandler:
		logger().Debugf("Option Handler: %s", h.Option())
		nego.optionHandlers[h.Option()] = h
	default:
		logger().Debugf("UNKNOWN HANDLER TYPE: %t", h)
	}
}

// Consider how to deal with or to ignore the packet
func (nego *Negotiator) Consider(p packet.Packet) {
	nego.chInput <- p
}

// Run the negotiator in a goroutine for processing any input and output packets
func (nego *Negotiator) Run(ctx context.Context) {
	for {
		select {
		case input := <-nego.chInput:
			nego.handle(ctx, input)
		case <-ctx.Done():
			nego.dispose()
			logger().Info("nego stopped")
			return
		}
	}
}

func (nego *Negotiator) dispose() {
	nego.BaseDispose()
}

func (nego *Negotiator) findOptionHandler(option telbyte.Option) OptionHandler {
	handler, ok := nego.optionHandlers[option]
	if !ok {
		logger().Warnf("[WARN] no handler: %s", option)
		return nil
	}
	return handler
}

func (nego *Negotiator) handle(ctx context.Context, input packet.Packet) {
	logger().Debugf("\x1b[32m<RECV>\x1b[0m %s\n", input)
	switch p := input.(type) {
	case *packet.CommandPacket:
		handler := nego.findOptionHandler(p.Option)
		if handler != nil {
			logger().Debugf("handshake on <%s %s>", p.Command, p.Option)
			go handler.Handshake(newOptionContext(ctx, handler, nego), p.Command)
		}
	case *packet.SubPacket:
		handler := nego.findOptionHandler(p.Option)
		if handler != nil {
			go handler.Subnegotiate(newOptionContext(ctx, handler, nego), p.Parameter)
		}
	}
}
