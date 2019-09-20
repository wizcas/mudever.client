package telnet

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/wizcas/mudever.svc/telnet/packet"
	"github.com/wizcas/mudever.svc/telnet/receiver"
	"github.com/wizcas/mudever.svc/telnet/stream"
)

type termErrorKind byte

const (
	termErrorSys termErrorKind = iota
	termErrorRecv
	termErrorSend
)

// Terminal is where the telnet process runs
type Terminal struct {
	Encoding TermEncoding
}

type terminalError struct {
	kind termErrorKind
	err  error
}

func (te terminalError) Error() string {
	var errName string
	switch te.kind {
	case termErrorSys:
		errName = "SYS"
	case termErrorRecv:
		errName = "RECV"
	case termErrorSend:
		errName = "SEND"
	default:
		errName = "UNKNOWN"
	}
	return fmt.Sprintf("[%s ERR] %s", errName, te.err.Error())
}

// NewTerminal creates a telnet terminal with specified encoding charset
func NewTerminal(encoding TermEncoding) *Terminal {
	return &Terminal{encoding}
}

func (t *Terminal) proc(r *stream.Reader, w *stream.Writer) error {
	chSendErr := make(chan error)
	recv := receiver.New(r)
	go recv.Run()
	go t.send(w, chSendErr)
	for {
		select {
		case pkt := <-recv.ChPacket:
			switch p := pkt.(type) {
			case *packet.DataPacket:
				output, err := t.decode(p.Data)
				if err != nil {
					return terminalError{termErrorSys, err}
				}
				os.Stdout.Write(output)
			case *packet.CommandPacket:
			case *packet.SubPacket:
				os.Stdout.Write([]byte(p.String()))
			}

		case err := <-recv.ChErr:
			return terminalError{termErrorRecv, err}
		case err := <-chSendErr:
			return terminalError{termErrorSend, err}
		}
	}
}

func (t *Terminal) decode(b []byte) ([]byte, error) {
	if t.Encoding == nil {
		return b, nil
	}
	return t.Encoding.NewDecoder().Bytes(b)
}

func (t *Terminal) encode(b []byte) ([]byte, error) {
	if t.Encoding == nil {
		return b, nil
	}
	return t.Encoding.NewEncoder().Bytes(b)
}

func (t *Terminal) send(w *stream.Writer, chErr chan error) {
	scanner := bufio.NewScanner(os.Stdin)
	buf := bytes.NewBuffer(nil)
	scanner.Split(splitLines)
	for scanner.Scan() {
		buf.Reset()
		buf.Write(scanner.Bytes())
		buf.Write(crlf)
		line, err := t.encode(buf.Bytes())
		if err != nil {
			chErr <- err
		}
		n, err := w.Write(line)
		if expect := len(line); n != expect {
			chErr <- fmt.Errorf("outgoing data loss: %d of %d bytes ar sent", n, expect)
		}
	}
}

func splitLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF {
		return 0, nil, nil
	}

	return bufio.ScanLines(data, atEOF)
}
