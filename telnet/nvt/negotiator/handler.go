package negotiator

import (
	"errors"

	"github.com/wizcas/mudever.svc/telnet/protocol"
)

// OptionHandler defines the common behavior of those who can process option negotiations
type OptionHandler interface {
	Option() protocol.OptByte
	Handshake(inCmd protocol.CmdByte) (protocol.CmdByte, error)
	Subnegotiate(inParameter []byte) ([]byte, error)
}

// Errors caused by handlers
var (
	ErrIgnore   = errors.New("IGNORE")
	ErrLackData = errors.New("LACK OF DATA")
)
