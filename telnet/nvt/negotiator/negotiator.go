package nego

import (
	"context"
	"log"

	"github.com/wizcas/mudever.svc/telnet/nvt/common"
	"github.com/wizcas/mudever.svc/telnet/packet"
	"github.com/wizcas/mudever.svc/telnet/telbyte"
)

// Negotiator takes care of telnet negotiations
type Negotiator struct {
	*common.SubProc
	controlHandlers map[telbyte.Command]ControlHandler
	optionHandlers  map[telbyte.Option]OptionHandler

	chInput      chan packet.Packet
	ChHandledCmd chan HandledCommand
	ChHandledSub chan HandledSub

	replier *replier
}

// New negotiator with an empty knowledge base
func New(chOutput chan<- packet.Packet) *Negotiator {
	return &Negotiator{
		SubProc:         common.NewSubProc(),
		optionHandlers:  make(map[telbyte.Option]OptionHandler),
		controlHandlers: make(map[telbyte.Command]ControlHandler),

		chInput:      make(chan packet.Packet),
		ChHandledCmd: make(chan HandledCommand),
		ChHandledSub: make(chan HandledSub),

		replier: newReplier(chOutput),
	}
}

// Know an handler makes the negotiator being able to deal with the corresponding incoming option
func (nego *Negotiator) Know(handler Handler) {
	switch h := handler.(type) {
	case ControlHandler:
		log.Printf("Control Handler: %s", h.Command())
		nego.controlHandlers[h.Command()] = h
	case OptionHandler:
		log.Printf("Option Handler: %s", h.Option())
		nego.optionHandlers[h.Option()] = h
	default:
		log.Println("UNKNOWN HANDLER TYPE")
	}
	handler.Initialize(nego)
}

// Consider how to deal with or to ignore the packet
func (nego *Negotiator) Consider(p packet.Packet) {
	nego.chInput <- p
}

// Run the negotiator in a goroutine for processing any input and output packets
func (nego *Negotiator) Run(ctx context.Context) {

	// Send queued replies when possible but don't block the input flow
	ctxReply, cancelReply := context.WithCancel(ctx)
	defer cancelReply()
	go nego.replier.run(ctxReply)

	for {
		select {
		case input := <-nego.chInput:
			nego.handle(input)
		case hcmd := <-nego.ChHandledCmd:
			p := packet.NewOptionCommandPacket(hcmd.Command, hcmd.Handler.Option())
			nego.replier.enqueue(p)
		case hsub := <-nego.ChHandledSub:
			p := packet.NewSubPacket(hsub.Handler.Option(), hsub.Parameter...)
			nego.replier.enqueue(p)
		case <-ctx.Done():
			nego.dispose()
			log.Println("nego stopped")
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
		log.Printf("[WARN] Handler for option <%s> not found", option)
		return nil
	}
	return handler
}

func (nego *Negotiator) handle(input packet.Packet) {
	log.Printf("\x1b[32m<NEGO RECV>\x1b[0m %s\n", input)
	switch p := input.(type) {
	case *packet.CommandPacket:
		handler := nego.findOptionHandler(p.Option)
		if handler != nil {
			log.Printf("handshake on <%s %s>", p.Command, p.Option)
			go handler.Handshake(p.Command)
		}
	case *packet.SubPacket:
		handler := nego.findOptionHandler(p.Option)
		if handler != nil {
			go handler.Subnegotiate(p.Parameter)
		}
	}
}
