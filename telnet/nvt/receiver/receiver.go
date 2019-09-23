package receiver

import (
	"bytes"
	"context"
	"io"

	"github.com/wizcas/mudever.svc/telnet/nvt/common"
	"github.com/wizcas/mudever.svc/telnet/packet"
	"go.uber.org/zap"
)

// Receiver reads data from network stream and parses it into Packets,
// which is comprehensible to Terminal
type Receiver struct {
	*common.SubProc
	src      io.Reader
	alloc    []byte
	chOutput chan packet.Packet

	processor *processor
	logger    *zap.SugaredLogger
}

// New telnet data receiver
func New(src io.Reader) *Receiver {
	return &Receiver{
		SubProc:   common.NewSubProc(),
		src:       src,
		alloc:     make([]byte, 128),
		chOutput:  make(chan packet.Packet),
		processor: newProcessor(),
	}
}

// Run the receiver loop to read and process data from incoming stream.
// Any processed packets are input into Receiver.ChPacket, and errors
// into Receiver.ChErr
func (r *Receiver) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			r.dispose()
			r.logger.Info("receiver stopped.")
			return
		default:
			data, err := r.readStream()
			go r.processor.proc(data, r.chOutput, r.GotError)
			if err != nil {
				r.GotError(err)
			}
		}
	}
}

// Output packet processed by received
func (r *Receiver) Output() <-chan packet.Packet {
	return r.chOutput
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
	close(r.chOutput)
	r.BaseDispose()
}
