package mtts

import (
	"fmt"

	"github.com/wizcas/mudever.svc/nvt/nego"
	"github.com/wizcas/mudever.svc/telbyte"
)

// MTTS Sub-commands
const (
	IS   = byte(0)
	SEND = byte(1)
)

// MUD Client Terminal Types
const (
	TypeDUMB  = "DUMB"
	TypeANSI  = "ANSI"
	TypeVT100 = "VT100"
	TypeXTERM = "XTERM"
)

// MTTS contains information for TerminalType negotiations.
type MTTS struct {
	// ClientName of the terminal including version preferably
	ClientName string
	// TerminalType which should be set to one of the 'Type*' enum values
	TerminalType string
	// Features indicates the terminal's capabilities and is the SUM of
	// the enabled features.
	Features   *featureSet
	queryTimes int
}

// New MTTS handler with default values.
func New(isUTF8 bool) *MTTS {
	feats := newFeatureSet(FeatANSI, Feat256Colors, FeatTrueColor, FeatMouseTracking)
	if isUTF8 {
		feats.add(FeatUTF8)
	}
	return &MTTS{
		ClientName:   "MUDEVER 0.1",
		TerminalType: TypeXTERM,
		Features:     feats,
		queryTimes:   0,
	}
}

// Option implements OptionHandler
func (h *MTTS) Option() telbyte.Option {
	return telbyte.TTYPE
}

// Handshake implements OptionHandler, it responds to only DO & DONT.
// Other commands will be ignored.
func (h *MTTS) Handshake(ctx *nego.OptionContext, inCmd telbyte.Command) {
	var res telbyte.Command
	switch inCmd {
	case telbyte.DO:
		res = telbyte.WILL
	case telbyte.DONT:
		res = telbyte.WONT
		h.queryTimes = 0
	default:
		return
	}
	ctx.SendCmd(res)
}

// Subnegotiate implements OptionHandler, and works in the way described at:
// https://tintin.sourceforge.io/protocols/mtts/
func (h *MTTS) Subnegotiate(ctx *nego.OptionContext, inParameter []byte) {
	if len(inParameter) == 0 {
		ctx.GotError(nego.ErrLackData)
		return
	}
	action := inParameter[0]
	if action != SEND {
		return
	}
	var payload []byte
	switch h.queryTimes {
	case 0:
		payload = []byte(h.ClientName)
	case 1:
		payload = []byte(h.TerminalType)
	default:
		payload = []byte(fmt.Sprintf("MTTS %d", h.Features.Value()))
	}
	params := append([]byte{IS}, payload...)
	ctx.SendSub(params)
	h.queryTimes++
}

func (h *MTTS) setQueryTimesForTesting(times int) {
	h.queryTimes = times
}
