package telnet

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
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
	var errorKindName string
	switch te.kind {
	case termErrorSys:
		errorKindName = "SYS"
	case termErrorRecv:
		errorKindName = "RECV"
	case termErrorSend:
		errorKindName = "SEND"
	default:
		errorKindName = "UNKNOWN"
	}
	return fmt.Sprintf("[%s ERR] %s", errorKindName, te.err.Error())
}

// NewTerminal creates a telnet terminal with specified encoding charset
func NewTerminal(encoding TermEncoding) *Terminal {
	return &Terminal{encoding}
}

func (t *Terminal) proc(r *reader, w *writer) error {
	chRecvErr := make(chan error)
	chSendErr := make(chan error)
	go t.recv(r, chRecvErr)
	go t.send(w, chSendErr)
	for {
		select {
		case err := <-chRecvErr:
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

func (t *Terminal) recv(r *reader, chErr chan error) {
	buf := make([]byte, 1)
	packet := bytes.NewBuffer(nil)
	for {
		packet.Reset()
		for {
			n, err := r.read(buf)
			if err != nil {
				if err == io.EOF {
					// chErr <- errors.New("hiohoho")
					break
				} else {
					chErr <- err
				}
			}
			if n > 0 {
				packet.Write(buf[:n])
			}
		}
		data := packet.Bytes()
		log.Println(data)
		output, err := t.decode(data)
		if err != nil {
			chErr <- err
		}
		os.Stdout.Write(output)
	}
}

func (t *Terminal) send(w *writer, chErr chan error) {
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
		n, err := w.write(line)
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
