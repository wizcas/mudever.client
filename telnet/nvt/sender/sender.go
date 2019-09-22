package sender

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/wizcas/mudever.svc/telnet/nvt/common"
	"github.com/wizcas/mudever.svc/telnet/packet"
)

// Sender takes packets, serialize them and write into destination writer.
type Sender struct {
	*common.SubProc
	sync.Mutex
	dst     io.Writer
	chInput chan packet.Packet
	running bool
}

// New sender to write packets into dst writer.
func New(dst io.Writer) *Sender {
	return &Sender{
		SubProc: common.NewSubProc(),
		dst:     dst,
	}
}

// Run the sender loop for sending packets
func (s *Sender) Run(ctx context.Context) {
	s.chInput = make(chan packet.Packet, 100)
	s.running = true
	for {
		select {
		case p := <-s.chInput:
			log.Printf("[SEND] input: %s", p)
			s.doSend(p)
		case <-ctx.Done():
			s.dispose()
			log.Println("sender stopped.")
			return
		}
	}
}

// Send put packet into its sending queue if available
func (s *Sender) Send(p packet.Packet) error {
	s.Lock()
	defer s.Unlock()
	if !s.running {
		return errors.New("sender is not running")
	}
	s.chInput <- p
	return nil
}

func (s *Sender) dispose() {
	s.Lock()
	defer s.Unlock()
	close(s.chInput)
	s.BaseDispose()
	s.running = false
}

func (s *Sender) doSend(p packet.Packet) {
	if p == nil {
		return
	}
	if data, err := p.Serialize(); err != nil {
		log.Printf("\x1b[31m<SEND ERR>\x1b[0m %v\n", err)
		s.GotError(err)
	} else {
		if n, err := s.dst.Write(data); err != nil {
			s.GotError(err)
		} else if n != len(data) {
			s.GotError(fmt.Errorf("data inconsistency: %d written (%d intended)", n, len(data)))
		} else {
			log.Printf("\x1b[31m<SEND SUCC>\x1b[0m %s\n", p)
		}
	}
}
