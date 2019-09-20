package sender

import (
	"fmt"
	"io"
	"log"

	"github.com/wizcas/mudever.svc/telnet/packet"
)

// Sender takes packets, serialize them and write into destination writer.
type Sender struct {
	dst     io.Writer
	ChInput chan packet.Packet
	ChErr   chan error
	ChStop  chan struct{}
}

// New sender to write packets into dst writer.
func New(dst io.Writer) *Sender {
	return &Sender{
		dst:     dst,
		ChInput: make(chan packet.Packet),
		ChErr:   make(chan error),
		ChStop:  make(chan struct{}),
	}
}

// Run the sender loop for sending packets
func (s *Sender) Run() {
LOOP:
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
		case <-s.ChStop:
			break LOOP
		}
	}
	s.dispose()
	log.Println("sender stopped.")
}

// Stop the sender and release resources
func (s *Sender) Stop() {
	s.ChStop <- struct{}{}
}

func (s *Sender) dispose() {
	close(s.ChInput)
	close(s.ChErr)
	close(s.ChStop)
}
