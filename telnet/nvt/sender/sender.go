package sender

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/wizcas/mudever.svc/telnet/nvt/common"
	"github.com/wizcas/mudever.svc/telnet/packet"
)

// Sender takes packets, serialize them and write into destination writer.
type Sender struct {
	*common.SubProc
	dst     io.Writer
	chInput chan packet.Packet
}

// New sender to write packets into dst writer.
func New(dst io.Writer) *Sender {
	return &Sender{
		SubProc: common.NewSubProc(),
		dst:     dst,
		chInput: make(chan packet.Packet),
	}
}

// Run the sender loop for sending packets
func (s *Sender) Run(ctx context.Context) {
	for {
		select {
		case p := <-s.chInput:
			if p != nil {
				log.Printf("\x1b[31m<SEND GOT>\x1b[0m %s\n", p)
				if data, err := p.Serialize(); err != nil {
					log.Printf("\x1b[31m<SEND ERR>\x1b[0m %v\n", err)
					s.ChErr <- err
				} else {
					if n, err := s.dst.Write(data); err != nil {
						s.ChErr <- err
					} else if n != len(data) {
						s.ChErr <- fmt.Errorf("data inconsistency: %d written (%d intended)", n, len(data))
					}
				}
				log.Printf("\x1b[31m<SEND SUCC>\x1b[0m %s\n", p)
			}
		case <-ctx.Done():
			s.dispose()
			log.Println("sender stopped.")
			return
		}
	}
}

// Input channel for pushing packet into sender
func (s *Sender) Input() chan<- packet.Packet {
	return s.chInput
}

func (s *Sender) dispose() {
	close(s.chInput)
	s.BaseDispose()
}
