package naws

import (
	"encoding/binary"

	"github.com/wizcas/mudever.svc/telnet/nvt/nego"
	"github.com/wizcas/mudever.svc/telnet/telbyte"
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
func (h *NAWS) Option() telbyte.Option {
	return telbyte.NAWS
}

// Handshake implements OptionHandler
func (h *NAWS) Handshake(ctx *nego.OptionContext, inCmd telbyte.Command) {
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
	ctx.SendCmd(reply)
	if h.submitting {
		ctx.SendSub(
			encodeWSValue(h.Width),
			encodeWSValue(h.Height),
		)
	}
}

// Subnegotiate implements OptionHandler
func (h *NAWS) Subnegotiate(ctx *nego.OptionContext, inParameter []byte) {
}

func encodeWSValue(val uint16) []byte {
	b2 := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(b2, val)
	return b2
}
