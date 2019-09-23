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

// TerminalError is an error reported by terminal
type TerminalError struct {
	kind  errorKind
	err   error
	panic bool
}

func newTerminalError(kind errorKind, err error, panic bool) TerminalError {
	return TerminalError{kind, err, panic}
}

func (te TerminalError) Panic() bool {
	return te.panic
}

func (te TerminalError) RawErr() error {
	return te.err
}

func (te TerminalError) Error() string {
	return fmt.Sprintf("[TERM %s ERR] %s", te.kind, te.err.Error())
}
