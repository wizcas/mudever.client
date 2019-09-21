package naws

import (
	"encoding/binary"

	nego "github.com/wizcas/mudever.svc/telnet/nvt/negotiator"
	"github.com/wizcas/mudever.svc/telnet/telbyte"
)

// NAWS stands for Negotiate About Window Size, which is used for
// dealing with telnet option 31
type NAWS struct {
	*nego.OptionHandlerBase
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
	return &NAWS{
		OptionHandlerBase: nego.NewOptionHandlerBase(),
	}
}

// Option implements OptionHandler
func (h *NAWS) Option() telbyte.Option {
	return telbyte.NAWS
}

// Handshake implements OptionHandler
func (h *NAWS) Handshake(inCmd telbyte.Command) {
	var reply telbyte.Command
	switch inCmd {
	case telbyte.DO:
		h.submitting = true
		reply = telbyte.WILL
	case telbyte.DONT:
		h.submitting = false
		reply = telbyte.WONT
	default:
		return
	}
	h.ChOutCmd <- nego.NewHandledCmd(h, reply)
	if h.submitting {
		h.ChOutSub <- nego.NewHandledSub(h, h.encodeParameter())
	}
}

// Subnegotiate implements OptionHandler
func (h *NAWS) Subnegotiate(inParameter []byte) {
}

func (h *NAWS) encodeParameter() []byte {
	bw := make([]byte, 2, 2)
	bh := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(bw, h.Width)
	binary.BigEndian.PutUint16(bh, h.Height)
	return append(bw, bh...)
}
