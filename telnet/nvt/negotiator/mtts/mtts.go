package mtts

import (
	"fmt"

	"github.com/wizcas/mudever.svc/telnet/nvt/negotiator"
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
		ClientName:   "MUDEVER 0.1",
		TerminalType: TypeXTERM,
		SupportFlag:  supportFlag,
		queryTimes:   0,
	}
}

// Option implements OptionHandler
func (h *MTTS) Option() telbyte.Option {
	return telbyte.TTYPE
}

// Handshake implements OptionHandler, it respond to only DO & DONT
func (h *MTTS) Handshake(inCmd telbyte.Command) (telbyte.Command, error) {
	switch inCmd {
	case telbyte.DO:
		return telbyte.WILL, nil
	case telbyte.DONT:
		h.queryTimes = 0
		return telbyte.WONT, nil
	default:
		return 0, negotiator.ErrIgnore
	}
}

// Subnegotiate implements OptionHandler, and works in the way described at:
// https://tintin.sourceforge.io/protocols/mtts/
func (h *MTTS) Subnegotiate(inParameter []byte) ([]byte, error) {
	if len(inParameter) == 0 {
		return nil, negotiator.ErrLackData
	}
	action := inParameter[0]
	if action != SEND {
		return nil, negotiator.ErrIgnore
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
	return append([]byte{IS}, payload...), nil
}
