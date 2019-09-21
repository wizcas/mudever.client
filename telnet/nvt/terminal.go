package nvt

import (
	"bufio"
	"bytes"
	"log"
	"os"

	nego "github.com/wizcas/mudever.svc/telnet/nvt/negotiator"
	"github.com/wizcas/mudever.svc/telnet/nvt/negotiator/mtts"
	"github.com/wizcas/mudever.svc/telnet/nvt/receiver"
	"github.com/wizcas/mudever.svc/telnet/nvt/sender"
	"github.com/wizcas/mudever.svc/telnet/packet"
	"github.com/wizcas/mudever.svc/telnet/stream"
)

// Terminal is where the telnet process runs
type Terminal struct {
	Encoding Encoding
	receiver *receiver.Receiver
	sender   *sender.Sender
}

// NewTerminal creates a telnet terminal with specified encoding charset
func NewTerminal(encoding Encoding) *Terminal {
	return &Terminal{
		Encoding: encoding,
	}
}

// Start the terminal session
func (t *Terminal) Start(r *stream.Reader, w *stream.Writer) error {

	ng := nego.New()
	ng.Know(mtts.New(false))
	go ng.Run()

	chInputErr := make(chan error)
	t.receiver = receiver.New(r)
	t.sender = sender.New(w)
	go t.receiver.Run()
	go t.sender.Run()
	go t.input(chInputErr)

	for {
		select {
		case pkt := <-t.receiver.ChOutput:
			switch p := pkt.(type) {
			case *packet.DataPacket:
				output, err := t.decode(p.Data)
				if err != nil {
					return terminalError{errorSys, err}
				}
				os.Stdout.Write(output)
			case *packet.CommandPacket, *packet.SubPacket:
				log.Printf("\x1b[32m<RECV>\x1b[0m %s\n", p)
				ng.ChInput <- p
				// res, err := nego.Handle(p)
				// if err != nil {
				// 	return err
				// }
				// if res != nil {
				// 	log.Printf("\x1b[36m<RPLY>\x1b[0m %s\n", res)
				// 	t.sender.ChInput <- res
				// }
			}
		case reply := <-ng.ChOutput:
			log.Printf("\x1b[36m<RPLY>\x1b[0m %s\n", reply)
			t.sender.ChInput <- reply
		case err := <-t.receiver.ChErr:
			return terminalError{errorRecv, err}
		case err := <-t.sender.ChErr:
			return terminalError{errorSend, err}
		case err := <-chInputErr:
			return terminalError{errorInput, err}
		case err := <-ng.ChErr:
			if err != nil && err != nego.ErrIgnore {
				return terminalError{errorNegotiator, err}
			}
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
		buf.Write(CRLF)
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
