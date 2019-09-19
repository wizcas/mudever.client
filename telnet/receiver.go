package telnet

import (
	"bytes"
	"log"
)

type receiver struct {
	src      *reader
	buffer   []byte
	chErr    chan error
	chPacket chan []byte
}

func newReceiver(src *reader) *receiver {
	return &receiver{
		src:      src,
		buffer:   make([]byte, 128),
		chErr:    make(chan error),
		chPacket: make(chan []byte),
	}
}

func (r *receiver) run() {
	for {
		if data, err := r.readStream(); err != nil {
			r.chErr <- err
		} else {
			// log.Println(data)
			log.Printf("[PACKET READ] len: %d", len(data))
			r.chPacket <- data
		}
	}
}

func (r *receiver) readStream() ([]byte, error) {
	packet := bytes.NewBuffer(nil)
	for {
		n, err := r.src.read(r.buffer)
		if n > 0 {
			packet.Write(r.buffer[:n])
		}
		if err != nil {
			if err == ErrEOS {
				break
			} else {
				return packet.Bytes(), err
			}
		}
	}
	return packet.Bytes(), nil
}
