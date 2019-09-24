package receiver

import (
	"bytes"
	"context"
	"io"
	"sync"

	"github.com/wizcas/mudever.svc/nvt/common"
	"github.com/wizcas/mudever.svc/packet"
)

// Receiver reads data from network stream and parses it into Packets,
// which is comprehensible to Terminal
type Receiver struct {
	sync.Mutex
	*common.BaseSubProc
	src      io.Reader
	alloc    []byte
	chRead   chan []byte
	chOutput chan packet.Packet
	running  bool

	processor *processor
}

// New telnet data receiver
func New(src io.Reader) *Receiver {
	return &Receiver{
		BaseSubProc: common.NewBaseSubProc(),
		src:         src,
		alloc:       make([]byte, 128),
		chRead:      make(chan []byte, 100),
		chOutput:    make(chan packet.Packet, 5),
		processor:   newProcessor(),
	}
}

// Run the receiver loop to read and process data from incoming stream.
// Any processed packets are input into Receiver.ChPacket, and errors
// into Receiver.ChErr
func (r *Receiver) Run(ctx context.Context) {
	r.setRunning(true)
	go r.feedData()
	for {
		select {
		case <-ctx.Done():
			r.dispose()
			common.Logger().Info("receiver stopped")
			return
		case data := <-r.chRead:
			go r.processor.proc(data, r.chOutput, r.GotError)
		}
	}
}

// Output packet processed by received
func (r *Receiver) Output() <-chan packet.Packet {
	return r.chOutput
}

func (r *Receiver) feedData() {
	for r.getRunning() {
		data, err := r.readStream()
		if !r.getRunning() {
			return
		}
		if err != nil {
			r.GotError(err)
		} else {
			r.chRead <- data
		}
	}
}

func (r *Receiver) setRunning(running bool) {
	r.Lock()
	defer r.Unlock()
	r.running = running
}
func (r *Receiver) getRunning() bool {
	r.Lock()
	defer r.Unlock()
	return r.running
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

func (r *Receiver) dispose() {
	r.setRunning(false)
	close(r.chRead)
	close(r.chOutput)
	r.BaseDispose()
}
