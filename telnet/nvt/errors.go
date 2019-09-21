package nvt

import "fmt"

type errorKind byte

const (
	errorSys errorKind = iota
	errorRecv
	errorSend
	errorInput
	errorNegotiator
)

var errorNames = map[errorKind]string{
	errorSys:        "SYS",
	errorRecv:       "RECV",
	errorSend:       "SEND",
	errorInput:      "INPUT",
	errorNegotiator: "NEGO",
}

func (k errorKind) String() string {
	name, ok := errorNames[k]
	if !ok {
		return "UNKNOWN"
	}
	return name
}

type terminalError struct {
	kind errorKind
	err  error
}

func (te terminalError) Error() string {
	return fmt.Sprintf("[TERM %s ERR] %s", te.kind, te.err.Error())
}
