package negotiator

import (
	"log"

	"github.com/wizcas/mudever.svc/telnet/packet"
	"github.com/wizcas/mudever.svc/telnet/protocol"
)

// Negotiator takes care of telnet negotiations
type Negotiator struct {
	controlHandlers map[protocol.CmdByte]ControlHandler
	optionHandlers  map[protocol.OptByte]OptionHandler
}

// New negotiator with an empty knowledge base
func New() *Negotiator {
	return &Negotiator{
		optionHandlers: make(map[protocol.OptByte]OptionHandler),
	}
}

// Know an handler makes the negotiator being able to deal with the corresponding incoming option
func (m *Negotiator) Know(handler Handler) {
	switch h := handler.(type) {
	case ControlHandler:
		log.Printf("Control Handler: %s", h.Command())
		m.controlHandlers[h.Command()] = h
	case OptionHandler:
		log.Printf("Option Handler: %s", h.Option())
		m.optionHandlers[h.Option()] = h
	default:
		log.Println("UNKNOWN HANDLER TYPE")
	}
}

func (m *Negotiator) findOptionHandler(option protocol.OptByte) OptionHandler {
	handler, ok := m.optionHandlers[option]
	if !ok {
		return nil
	}
	return handler
}

// Handle an incoming packet by the option within, and returns an outgoing packet.
// A nil packet will be returned if no handler found, income ignored, or error occured.
// If there's an error needs to be dealt with, it'll be set in the second return value.
func (m *Negotiator) Handle(p packet.Packet) (packet.Packet, error) {
	if cmd, ok := p.(*packet.CommandPacket); ok {
		handler := m.findOptionHandler(cmd.Option)
		if handler == nil {
			return nil, nil
		}
		out, err := handler.Handshake(cmd.Command)
		if err != nil {
			return nil, filterError(err)
		}
		return packet.NewOptionCommandPacket(out, cmd.Option), nil
	}

	if sub, ok := p.(*packet.SubPacket); ok {
		handler := m.findOptionHandler(sub.Option)
		if handler == nil {
			return nil, nil
		}
		out, err := handler.Subnegotiate(sub.Parameter)
		if err != nil {
			return nil, filterError(err)
		}
		return packet.NewSubPacket(sub.Option, out), nil
	}

	return nil, nil
}

func filterError(err error) error {
	if err == ErrIgnore {
		return nil
	}
	return err
}
