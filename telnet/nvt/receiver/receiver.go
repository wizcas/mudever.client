package receiver

import (
	"bytes"
	"io"
	"log"

	"github.com/wizcas/mudever.svc/telnet/packet"
)

// Receiver reads data from network stream and parses it into Packets,
// which is comprehensible to Terminal
type Receiver struct {
	src      io.Reader
	alloc    []byte
	ChErr    chan error
	ChOutput chan packet.Packet

	processor *processor
}

// New telnet data receiver
func New(src io.Reader) *Receiver {
	return &Receiver{
		src:       src,
		alloc:     make([]byte, 128),
		ChErr:     make(chan error),
		ChOutput:  make(chan packet.Packet),
		processor: newProcessor(),
	}
}

// Run the receiver loop to read and process data from incoming stream.
// Any processed packets are input into Receiver.ChPacket, and errors
// into Receiver.ChErr
func (r *Receiver) Run() {
	for {
		data, err := r.readStream()
		log.Printf("[PACKET RECV] len: %d", len(data))
		r.processor.proc(data, r.ChOutput, r.ChErr)
		if err != nil {
			r.ChErr <- err
		}
	}
}

func (r *Receiver) readStream() ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	for {
		n, err := r.src.Read(r.alloc)
		if n > 0 {
			buffer.Write(r.alloc[:n])
		}
		if err != nil {
			return buffer.Bytes(), err
		}
		if n == 0 { // (0, nil) indicates the end of stream
			break
		}
	}
	return buffer.Bytes(), nil
}
