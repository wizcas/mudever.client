package telnet

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/wizcas/mudever.svc/telnet/packet"
	"github.com/wizcas/mudever.svc/telnet/receiver"
	"github.com/wizcas/mudever.svc/telnet/sender"
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
	receiver *receiver.Receiver
	sender   *sender.Sender
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
	return &Terminal{
		Encoding: encoding,
	}
}

// Start the terminal session
func (t *Terminal) Start(r *stream.Reader, w *stream.Writer) error {
	chSendErr := make(chan error)
	t.receiver = receiver.New(r)
	t.sender = sender.New(w)
	go t.receiver.Run()
	go t.sender.Run()
	go t.input(chSendErr)

	for {
		select {
		case pkt := <-t.receiver.ChPacket:
			switch p := pkt.(type) {
			case *packet.DataPacket:
				output, err := t.decode(p.Data)
				if err != nil {
					return terminalError{termErrorSys, err}
				}
				os.Stdout.Write(output)
			case *packet.CommandPacket, *packet.SubPacket:
				log.Println(p)
			}
		case err := <-t.receiver.ChErr:
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

func (t *Terminal) input(chErr chan error) {
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
		linePacket := packet.NewDataPacket(line)
		t.sender.ChInput <- linePacket
	}
}

func splitLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF {
		return 0, nil, nil
	}

	return bufio.ScanLines(data, atEOF)
}
