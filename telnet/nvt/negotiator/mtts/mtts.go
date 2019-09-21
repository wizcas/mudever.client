package mtts

import (
	"fmt"
	"log"

	nego "github.com/wizcas/mudever.svc/telnet/nvt/negotiator"
	"github.com/wizcas/mudever.svc/telnet/telbyte"
)

// MTTS Sub-commands
const (
	IS   = 0
	SEND = 1
)

// MUD Client Terminal Types
const (
	TypeDUMB  = "DUMB"
	TypeANSI  = "ANSI"
	TypeVT100 = "VT100"
	TypeXTERM = "XTERM"
)

// MUD Terminal Type Bit Values
const (
	SupportANSI uint = 1 << iota
	SupportVT100
	SupportUTF8
	Support256Colors
	SupportMouseTracking
	SupportOscColorPalette
	SupportScreenReader
	SupportProxy
	SupportTrueColor
)

// MTTS contains information for TerminalType negotiations.
type MTTS struct {
	*nego.OptionHandlerBase
	// ClientName of the terminal including version preferably
	ClientName string
	// TerminalType which should be set to one of the 'Type*' enum values
	TerminalType string
	// SupportFlag indicates the terminal's capabilities,
	// which should be the SUM of all desired 'Support*' enum values
	SupportFlag uint
	queryTimes  int
}

// New MTTS handler with default values.
func New(isUTF8 bool) *MTTS {
	supportFlag := SupportANSI + Support256Colors + SupportTrueColor + SupportMouseTracking
	if isUTF8 {
		supportFlag += SupportUTF8
	}
	return &MTTS{
		OptionHandlerBase: nego.NewOptionHandlerBase(),
		ClientName:        "MUDEVER 0.1",
		TerminalType:      TypeXTERM,
		SupportFlag:       supportFlag,
		queryTimes:        0,
	}
}

// Option implements OptionHandler
func (h *MTTS) Option() telbyte.Option {
	return telbyte.TTYPE
}

// Handshake implements OptionHandler, it responds to only DO & DONT.
// Other commands will be ignored.
func (h *MTTS) Handshake(inCmd telbyte.Command) {
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
	h.ChOutCmd <- nego.NewHandledCmd(h, res)
	log.Printf("MTTS replys %s", res)
}

// Subnegotiate implements OptionHandler, and works in the way described at:
// https://tintin.sourceforge.io/protocols/mtts/
func (h *MTTS) Subnegotiate(inParameter []byte) {
	if len(inParameter) == 0 {
		h.ChErr <- nego.ErrLackData
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
		payload = []byte(fmt.Sprintf("MTTS %d", h.SupportFlag))
	}
	params := append([]byte{IS}, payload...)
	h.ChOutSub <- nego.NewHandledSub(h, params)
}
