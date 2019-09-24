package nego

import (
	"context"

	"github.com/logrusorgru/aurora"
	"github.com/wizcas/mudever.svc/nvt/common"
	"github.com/wizcas/mudever.svc/packet"
	"github.com/wizcas/mudever.svc/telbyte"
)

// Negotiator takes care of telnet negotiations
type Negotiator struct {
	*common.BaseSubProc
	controlHandlers map[telbyte.Command]ControlHandler
	optionHandlers  map[telbyte.Option]OptionHandler

	chInput chan packet.Packet
	sender  common.PacketSender
}

// New negotiator takes a PacketSender for its handlers to feed replies.
// It is created with an empty knowledge base, i.e. you'll have to feed in
// Handlers with Know() method.
func New(sender common.PacketSender) *Negotiator {
	return &Negotiator{
		BaseSubProc:     common.NewBaseSubProc(),
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
		common.Logger().Debug(aurora.Yellow("nego knows (control)"), h.Command())
		nego.controlHandlers[h.Command()] = h
	case OptionHandler:
		common.Logger().Debug(aurora.Yellow("nego knows (option)"), h.Option())
		nego.optionHandlers[h.Option()] = h
	default:
		common.Logger().Debug(aurora.Red("unsupported handler"), h)
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
			common.Logger().Info("nego stopped")
			return
		}
	}
}

// Commit a packet to the negotiator's sender
func (nego *Negotiator) Commit(p packet.Packet) error {
	return nego.sender.Send(p)
}

func (nego *Negotiator) dispose() {
	nego.BaseDispose()
}

func (nego *Negotiator) findControlHandler(cmd telbyte.Command) ControlHandler {
	handler, ok := nego.controlHandlers[cmd]
	if !ok {
		common.Logger().Debug(aurora.Red("no CONTROL handler for"), cmd)
		return nil
	}
	return handler
}

func (nego *Negotiator) findOptionHandler(option telbyte.Option) OptionHandler {
	handler, ok := nego.optionHandlers[option]
	if !ok {
		common.Logger().Debug(aurora.Red("no OPTION handler for"), option)
		return nil
	}
	return handler
}

func (nego *Negotiator) handle(ctx context.Context, input packet.Packet) {
	common.Logger().Debug(aurora.Blue("negotiating"), input)
	switch p := input.(type) {
	case *packet.CommandPacket:
		if p.IsOption() {
			if handler := nego.findOptionHandler(p.Option); handler != nil {
				go handler.Handshake(NewOptionContext(ctx, handler, nego), p.Command)
			}
		} else {
			if handler := nego.findControlHandler(p.Command); handler != nil {
				// TODO: call control functions
			}
		}
	case *packet.SubPacket:
		handler := nego.findOptionHandler(p.Option)
		if handler != nil {
			go handler.Subnegotiate(NewOptionContext(ctx, handler, nego), p.Parameter)
		}
	}
}
