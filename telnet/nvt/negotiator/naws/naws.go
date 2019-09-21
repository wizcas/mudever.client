package naws

import (
	"github.com/wizcas/mudever.svc/telnet/nvt/negotiator"
	"github.com/wizcas/mudever.svc/telnet/protocol"
)

// NAWS stands for Negotiate About Window Size, which is used for
// dealing with telnet option 31
type NAWS struct {
	// Width is defined as how many characters can be displayed horizontally.
	// Value 0 means let server side decide the width.
	Width uint16
	// Height is defined as how many lines can be displayed vertically.
	// Value 0 means let the server decide the width.
	Height uint16

	submitting bool
}

// New an NAWS object
func New() *NAWS {
	return &NAWS{}
}

// Option implements OptionHandler
func (h *NAWS) Option() protocol.OptByte {
	return protocol.NAWS
}

// Handshake implements OptionHandler
func (h *NAWS) Handshake(inCmd protocol.CmdByte) (protocol.CmdByte, error) {
	switch inCmd {
	case protocol.DO:
		h.submitting = true
		return protocol.WILL, nil
	case protocol.DONT:
		h.submitting = false
		return protocol.WONT, nil
	default:
		return 0, negotiator.ErrIgnore
	}
}

// Subnegotiate implements OptionHandler
func (h *NAWS) Subnegotiate(inParameter []byte) ([]byte, error) {
	return nil, negotiator.ErrIgnore
}
