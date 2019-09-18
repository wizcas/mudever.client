package telnet

import (
	"bufio"
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

func (block *readBlock) eof() bool {
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

func newReader(r io.Reader) *reader {
	return &reader{
		source:   r,
		buffered: bufio.NewReader(r),
		inPacket: false,
		traffic:  0,
	}
}

func (r *reader) packetEnds() bool {
	return r.buffered.Buffered() <= 0 && r.inPacket
}

func (r *reader) read(data []byte) (int, error) {
	block := newReadBlock(data)
	for !block.eof() {
		if r.packetEnds() {
			r.inPacket = false
			return block.size, io.EOF
		}
		b, err := r.buffered.ReadByte()
		r.inPacket = true
		if err != nil {
			return block.size, err
		}
		if b == 255 {
			next, err := r.buffered.Peek(1)
			if err != nil {
				return block.size, err
			}
			log.Printf("[IAC] %d\n", next[0])
		}
		block.writeByte(b)
		r.traffic++
	}
	return block.size, nil
}
