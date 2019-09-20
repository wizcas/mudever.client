package receiver

import (
	"bytes"
	"io"
	"log"

	"github.com/wizcas/mudever.svc/telnet/packet"
	"github.com/wizcas/mudever.svc/telnet/stream"
)

// Receiver reads data from network stream and parses it into Packets,
// which is comprehensible to Terminal
type Receiver struct {
	src      io.Reader
	alloc    []byte
	ChErr    chan error
	ChPacket chan packet.Packet

	state *processor
}

// New telnet data receiver
func New(src io.Reader) *Receiver {
	return &Receiver{
		src:      src,
		alloc:    make([]byte, 128),
		ChErr:    make(chan error),
		ChPacket: make(chan packet.Packet),
		state:    newProcessor(),
	}
}

// Run the receiver loop to read and process data from incoming stream.
// Any processed packets are input into Receiver.ChPacket, and errors
// into Receiver.ChErr
func (r *Receiver) Run() {
	for {
		data, err := r.readStream()
		if err != nil {
			r.ChErr <- err
		}
		log.Printf("[PACKET RECV] len: %d", len(data))
		r.state.proc(data, r.ChPacket, r.ChErr)
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
			if err == stream.ErrEOS {
				break
			} else {
				return buffer.Bytes(), err
			}
		}
	}
	return buffer.Bytes(), nil
}
