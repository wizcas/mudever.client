package telnet

import (
	"bufio"
	"errors"
	"io"
	"log"
)

type readBlock struct {
	buffer []byte
	wip    []byte
	size   int
}

func newReadBlock(buffer []byte) *readBlock {
	return &readBlock{buffer, buffer, 0}
}

func (block *readBlock) exhausted() bool {
	return len(block.wip) <= 0
}

func (block *readBlock) writeByte(b byte) {
	block.wip[0] = b
	block.wip = block.wip[1:]
	block.size++
}

type reader struct {
	source   io.Reader
	buffered *bufio.Reader
	inPacket bool
	traffic  uint64
}

// ErrEOS stands for 'End of Stream'.
// It means that currently the stream transmission is done,
// but connection is still alive for more incoming stream.
// The buffered data is considered complete when streaming stops, and
// The terminal should proceed the buffered data on EOS signal.
var ErrEOS = errors.New("END OF STREAM")

func newReader(r io.Reader) *reader {
	return &reader{
		source:   r,
		buffered: bufio.NewReader(r),
		inPacket: false,
		traffic:  0,
	}
}

func (r *reader) streamEnds() bool {
	return r.buffered.Buffered() <= 0 && r.inPacket
}

func (r *reader) read(data []byte) (int, error) {
	block := newReadBlock(data)
	for !block.exhausted() {
		if r.streamEnds() {
			r.inPacket = false
			log.Printf("[BLOCK READ] (EOS) len: %d", block.size)
			return block.size, ErrEOS
		}
		// Check for EOF in case the reader is closed
		if _, err := r.buffered.Peek(1); err != nil {
			log.Printf("[BLOCK READ] (EOF) len: %d", block.size)
			return block.size, err
		}

		b, err := r.buffered.ReadByte()
		r.inPacket = true
		if err != nil {
			return block.size, err
		}
		if b == 255 {
			next, err := r.buffered.Peek(1)
			if err != nil {
				log.Printf("[BLOCK READ] (IAC ERR) len: %d", block.size)
				return block.size, err
			}
			log.Printf("[IAC] %d\n", next[0])
		}
		block.writeByte(b)
		r.traffic++
	}
	log.Printf("[BLOCK READ] (BUF END) len: %d", block.size)
	return block.size, nil
}
