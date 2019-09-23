package nvt

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os"
	"sync"

	"github.com/wizcas/mudever.svc/telnet/nvt/common"
	"github.com/wizcas/mudever.svc/telnet/nvt/nego"
	"github.com/wizcas/mudever.svc/telnet/nvt/nego/mtts"
	"github.com/wizcas/mudever.svc/telnet/nvt/nego/naws"
	"github.com/wizcas/mudever.svc/telnet/nvt/receiver"
	"github.com/wizcas/mudever.svc/telnet/nvt/sender"
	"github.com/wizcas/mudever.svc/telnet/packet"
	"github.com/wizcas/mudever.svc/telnet/stream"
)

// Terminal is where the telnet process runs
type Terminal struct {
	Encoding  Encoding
	receiver  *receiver.Receiver
	sender    *sender.Sender
	chStopped chan struct{}
}

// NewTerminal creates a telnet terminal with specified encoding charset
func NewTerminal(encoding Encoding) *Terminal {
	return &Terminal{
		Encoding:  encoding,
		chStopped: make(chan struct{}),
	}
}

// Stopped returns a channel that emits value only when this terminal has stopped running
func (t *Terminal) Stopped() <-chan struct{} {
	return t.chStopped
}

// Start the terminal session
func (t *Terminal) Start(rootCtx context.Context, r *stream.Reader, w *stream.Writer, chErr chan TerminalError) {
	wg := sync.WaitGroup{}
	t.receiver = receiver.New(r)
	t.startSubProc(rootCtx, t.receiver, &wg)

	t.sender = sender.New(w)
	t.startSubProc(rootCtx, t.sender, &wg)

	chInputErr := make(chan error)
	go t.input(chInputErr)

	ng := nego.New(t.sender)
	ng.Know(mtts.New(false))
	ng.Know(naws.New())
	t.startSubProc(rootCtx, ng, &wg)
LOOP:
	for {
		select {
		case pkt := <-t.receiver.Output():
			switch p := pkt.(type) {
			case *packet.DataPacket:
				output, err := t.decode(p.Data)
				if err != nil {
					chErr <- newTerminalError(errorRecv, err, false)
				}
				os.Stdout.Write(output)
			case *packet.CommandPacket, *packet.SubPacket:
				// Send commands and subnegotiations to Negotiator
				ng.Consider(p)
			}
		case err := <-t.receiver.Err():
			chErr <- newTerminalError(errorRecv, err, err == io.EOF)
		case err := <-ng.Err():
			chErr <- newTerminalError(errorNegotiator, err, false)
		case err := <-t.sender.Err():
			chErr <- newTerminalError(errorSend, err, err == io.EOF)
		case err := <-chInputErr:
			chErr <- newTerminalError(errorInput, err, false)
		case <-rootCtx.Done():
			break LOOP
		}
	}

	wg.Wait()
	common.Logger().Info("terminal stopped")
	close(t.chStopped)
}

func (t *Terminal) startSubProc(ctx context.Context, subproc common.SubProc, wg *sync.WaitGroup) {
	wg.Add(1)
	go func(sp common.SubProc) {
		childCtx, _ := context.WithCancel(ctx)
		go subproc.Run(childCtx)
		select {
		case <-sp.Stopped():
			wg.Done()
		}
	}(subproc)
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
