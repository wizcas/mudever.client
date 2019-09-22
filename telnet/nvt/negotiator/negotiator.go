package nego

import (
	"container/list"
	"context"
	"log"
	"sync"

	"github.com/wizcas/mudever.svc/telnet/nvt/common"
	"github.com/wizcas/mudever.svc/telnet/packet"
	"github.com/wizcas/mudever.svc/telnet/telbyte"
)

// Negotiator takes care of telnet negotiations
type Negotiator struct {
	*common.SubProc
	controlHandlers map[telbyte.Command]ControlHandler
	optionHandlers  map[telbyte.Option]OptionHandler

	chOutput     chan<- packet.Packet
	chInput      chan packet.Packet
	ChHandledCmd chan HandledCommand
	ChHandledSub chan HandledSub

	replyQ    *list.List
	replySync *sync.Mutex
}

// New negotiator with an empty knowledge base
func New(chOutput chan<- packet.Packet) *Negotiator {
	return &Negotiator{
		SubProc:         common.NewSubProc(),
		optionHandlers:  make(map[telbyte.Option]OptionHandler),
		controlHandlers: make(map[telbyte.Command]ControlHandler),

		chOutput:     chOutput,
		chInput:      make(chan packet.Packet, 5),
		ChHandledCmd: make(chan HandledCommand, 5),
		ChHandledSub: make(chan HandledSub, 5),

		replyQ:    list.New(),
		replySync: &sync.Mutex{},
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
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				p := nego.deReply()
				if p != nil {
					nego.chOutput <- p
				}
			}
		}
	}(ctxReply)

LOOP:
	for {
		select {
		case input := <-nego.chInput:
			nego.handle(input)
		case hcmd := <-nego.ChHandledCmd:
			p := packet.NewOptionCommandPacket(hcmd.Command, hcmd.Handler.Option())
			nego.enReply(p)
		case hsub := <-nego.ChHandledSub:
			p := packet.NewSubPacket(hsub.Handler.Option(), hsub.Parameter...)
			nego.enReply(p)
		case <-ctx.Done():
			break LOOP
		}
	}
	nego.dispose()
	log.Println("nego stopped")
}

func (nego *Negotiator) dispose() {
	nego.replyQ.Init()
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

func (nego *Negotiator) enReply(p packet.Packet) {
	nego.replySync.Lock()
	defer nego.replySync.Unlock()
	nego.replyQ.PushBack(p)
}

func (nego *Negotiator) deReply() packet.Packet {
	nego.replySync.Lock()
	e := nego.replyQ.Front()
	if e == nil {
		nego.replySync.Unlock()
		return nil
	}
	nego.replyQ.Remove(e)
	nego.replySync.Unlock()
	reply, ok := e.Value.(packet.Packet)
	if !ok {
		log.Printf("error element in nego outgoing queue: %v", e.Value)
		return nil
	}
	log.Printf("\x1b[36m<NEGO RPLY>\x1b[0m %s\n", reply)
	return reply
}
