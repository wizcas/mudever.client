package nego

import (
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

	ChInput      chan packet.Packet
	ChOutput     chan packet.Packet
	ChHandledCmd chan HandledCommand
	ChHandledSub chan HandledSub
}

// New negotiator with an empty knowledge base
func New() *Negotiator {
	return &Negotiator{
		SubProc:         common.NewSubProc(),
		optionHandlers:  make(map[telbyte.Option]OptionHandler),
		controlHandlers: make(map[telbyte.Command]ControlHandler),

		ChInput:      make(chan packet.Packet, 5),
		ChOutput:     make(chan packet.Packet, 5),
		ChHandledCmd: make(chan HandledCommand, 5),
		ChHandledSub: make(chan HandledSub, 5),
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

// Run the negotiator in a goroutine for processing any input and output packets
func (nego *Negotiator) Run() {
LOOP:
	for {
		select {
		case input := <-nego.ChInput:
			log.Printf("negotiating: %s", input)
			nego.handle(input)
		case hcmd := <-nego.ChHandledCmd:
			log.Printf("handled cmd: %s %s", hcmd.Command, hcmd.Handler.Option())
			nego.ChOutput <- packet.NewOptionCommandPacket(hcmd.Command, hcmd.Handler.Option())
		case hsub := <-nego.ChHandledSub:
			nego.ChOutput <- packet.NewSubPacket(hsub.Handler.Option(), hsub.Parameter...)
		case <-nego.ChStop:
			break LOOP
		}
	}
	nego.dispose()
	log.Println("nego stopped")
}

func (nego *Negotiator) dispose() {
	close(nego.ChInput)
	close(nego.ChOutput)
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

func (nego *Negotiator) handle(p packet.Packet) {
	if cmd, ok := p.(*packet.CommandPacket); ok {
		handler := nego.findOptionHandler(cmd.Option)
		if handler == nil {
			return
		}
		log.Printf("handshake on <%s %s>", cmd.Command, cmd.Option)
		go handler.Handshake(cmd.Command)
	} else if sub, ok := p.(*packet.SubPacket); ok {
		handler := nego.findOptionHandler(sub.Option)
		if handler == nil {
			return
		}
		go handler.Subnegotiate(sub.Parameter)
	}
}
