package sender

import (
	"fmt"
	"io"

	"github.com/wizcas/mudever.svc/telnet/packet"
)

// Sender takes packets, serialize them and write into destination writer.
type Sender struct {
	dst     io.Writer
	ChInput chan packet.Packet
	ChErr   chan error
}

// New sender to write packets into dst writer.
func New(dst io.Writer) *Sender {
	return &Sender{
		dst:     dst,
		ChInput: make(chan packet.Packet),
		ChErr:   make(chan error),
	}
}

// Run the sender loop for sending packets
func (s *Sender) Run() {
	for {
		select {
		case p := <-s.ChInput:
			if data, err := p.Serialize(); err != nil {
				s.ChErr <- err
			} else {
				if n, err := s.dst.Write(data); err != nil {
					s.ChErr <- err
				} else if n != len(data) {
					s.ChErr <- fmt.Errorf("data inconsistency: %d written (%d intended)", n, len(data))
				}
			}
		}
	}
}
