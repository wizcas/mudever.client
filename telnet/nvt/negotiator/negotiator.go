package nego

import (
	"container/list"
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
	chOutput     chan<- packet.Packet
	ChHandledCmd chan HandledCommand
	ChHandledSub chan HandledSub

	replyQ *list.List
}

// New negotiator with an empty knowledge base
func New(chOutput chan<- packet.Packet) *Negotiator {
	return &Negotiator{
		SubProc:         common.NewSubProc(),
		optionHandlers:  make(map[telbyte.Option]OptionHandler),
		controlHandlers: make(map[telbyte.Command]ControlHandler),

		chInput:      make(chan packet.Packet, 5),
		chOutput:     chOutput,
		ChHandledCmd: make(chan HandledCommand, 5),
		ChHandledSub: make(chan HandledSub, 5),

		replyQ: list.New(),
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
LOOP:
	for {
		select {
		case input := <-nego.chInput:
			log.Printf("negotiating: %s", input)
			nego.handle(input)
		case hcmd := <-nego.ChHandledCmd:
			log.Printf("handled cmd: %s %s", hcmd.Command, hcmd.Handler.Option())
			p := packet.NewOptionCommandPacket(hcmd.Command, hcmd.Handler.Option())
			nego.enqReply(p)
		case hsub := <-nego.ChHandledSub:
			p := packet.NewSubPacket(hsub.Handler.Option(), hsub.Parameter...)
			nego.enqReply(p)
		case nego.chOutput <- nego.deqReply():
		case <-ctx.Done():
			break LOOP
		}
	}
	nego.dispose()
	log.Println("nego stopped")
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
	if cmd, ok := input.(*packet.CommandPacket); ok {
		handler := nego.findOptionHandler(cmd.Option)
		if handler == nil {
			return
		}
		log.Printf("handshake on <%s %s>", cmd.Command, cmd.Option)
		go handler.Handshake(cmd.Command)
	} else if sub, ok := input.(*packet.SubPacket); ok {
		handler := nego.findOptionHandler(sub.Option)
		if handler == nil {
			return
		}
		go handler.Subnegotiate(sub.Parameter)
	}
}

func (nego *Negotiator) enqReply(p packet.Packet) {
	nego.replyQ.PushBack(p)
}

func (nego *Negotiator) deqReply() packet.Packet {
	e := nego.replyQ.Front()
	if e == nil {
		return nil
	}
	nego.replyQ.Remove(e)
	reply, ok := e.Value.(packet.Packet)
	if !ok {
		log.Printf("error element in nego outgoing queue: %v", e.Value)
		return nil
	}
	log.Printf("\x1b[36m<NEGO RPLY>\x1b[0m %s\n", reply)
	return reply
}
