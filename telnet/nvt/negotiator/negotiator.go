package negotiator

import (
	"github.com/wizcas/mudever.svc/telnet/packet"
	"github.com/wizcas/mudever.svc/telnet/protocol"
)

// Negotiator takes care of telnet negotiations
type Negotiator struct {
	OptionMap map[protocol.OptByte]OptionHandler
}

// New negotiator with an empty knowledge base
func New() *Negotiator {
	return &Negotiator{
		OptionMap: make(map[protocol.OptByte]OptionHandler),
	}
}

// Know an handler makes the negotiator being able to deal with the corresponding incoming option
func (m *Negotiator) Know(handler OptionHandler) {
	m.OptionMap[handler.Option()] = handler
}

func (m *Negotiator) findHandler(option protocol.OptByte) OptionHandler {
	handler, ok := m.OptionMap[option]
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
		handler := m.findHandler(cmd.Option)
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
		handler := m.findHandler(sub.Option)
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
