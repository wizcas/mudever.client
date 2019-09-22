package nvt

import (
	"bufio"
	"bytes"
	"context"
	"os"

	nego "github.com/wizcas/mudever.svc/telnet/nvt/negotiator"
	"github.com/wizcas/mudever.svc/telnet/nvt/negotiator/mtts"
	"github.com/wizcas/mudever.svc/telnet/nvt/negotiator/naws"
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
	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.receiver = receiver.New(log(), r)
	recvCtx, _ := context.WithCancel(rootCtx)
	go t.receiver.Run(recvCtx)
	t.sender = sender.New(w)
	sendCtx, _ := context.WithCancel(rootCtx)
	go t.sender.Run(sendCtx)
	chInputErr := make(chan error)
	go t.input(chInputErr)

	ng := nego.New(t.sender)
	ng.Know(mtts.New(false))
	ng.Know(naws.New())
	ngCtx, _ := context.WithCancel(rootCtx)
	go ng.Run(ngCtx)

	for {
		select {
		case pkt := <-t.receiver.Output():
			switch p := pkt.(type) {
			case *packet.DataPacket:
				output, err := t.decode(p.Data)
				if err != nil {
					return terminalError{errorSys, err}
				}
				os.Stdout.Write(output)
			case *packet.CommandPacket, *packet.SubPacket:
				// Send commands and subnegotiations to Negotiator
				ng.Consider(p)
			}
		case err := <-t.receiver.Err():
			return terminalError{errorRecv, err}
		case err := <-ng.Err():
			return terminalError{errorNegotiator, err}
		case err := <-t.sender.Err():
			return terminalError{errorSend, err}
		case err := <-chInputErr:
			return terminalError{errorInput, err}
		case <-rootCtx.Done():
			return rootCtx.Err()
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
		if err := t.sender.Send(linePacket); err != nil {
			chErr <- err
		}
	}
}

func splitLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF {
		return 0, nil, nil
	}

	return bufio.ScanLines(data, atEOF)
}
