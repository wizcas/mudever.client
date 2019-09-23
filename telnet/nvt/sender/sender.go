package sender

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/logrusorgru/aurora"
	"github.com/wizcas/mudever.svc/telnet/nvt/common"
	"github.com/wizcas/mudever.svc/telnet/packet"
)

// Sender takes packets, serialize them and write into destination writer.
type Sender struct {
	*common.BaseSubProc
	sync.Mutex
	dst     io.Writer
	chInput chan packet.Packet
	running bool
}

// New sender to write packets into dst writer.
func New(dst io.Writer) *Sender {
	return &Sender{
		BaseSubProc: common.NewBaseSubProc(),
		dst:         dst,
	}
}

// Run the sender loop for sending packets
func (s *Sender) Run(ctx context.Context) {
	s.chInput = make(chan packet.Packet, 100)
	s.running = true
	for {
		select {
		case p := <-s.chInput:
			common.Logger().Debug(aurora.Green("sending:"), p)
			s.doSend(p)
		case <-ctx.Done():
			s.dispose()
			common.Logger().Info("sender stopped")
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
		s.GotError(err)
	} else {
		if n, err := s.dst.Write(data); err != nil {
			s.GotError(err)
		} else if n != len(data) {
			s.GotError(fmt.Errorf("data inconsistency: %d written (%d intended)", n, len(data)))
		} else {
			common.Logger().Debug(aurora.Green("sent"), p)
		}
	}
}
