package stream

import (
	"bufio"
	"io"
)

type readBlock struct {
	buffer   []byte
	wip      []byte
	readSize int
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
	block.readSize++
}

// Reader reads bytes from a reader that keep alive
// (i.e. not necessarily throws EOF when there's no more data),
// but may stop feeding data from time to time.
// The reader detects pauses in data transimission. It indicates the
// end of a stream by the Read return value of (0, nil)
type Reader struct {
	source    io.Reader
	buffered  *bufio.Reader
	streaming bool
	wantEOS   bool
	traffic   uint64
}

// NewReader returns a new reader
func NewReader(r io.Reader) *Reader {
	return &Reader{
		source:    r,
		buffered:  bufio.NewReader(r),
		streaming: false,
		wantEOS:   false,
		traffic:   0,
	}
}

func (r *Reader) streamEnds() bool {
	return r.buffered.Buffered() <= 0 && r.streaming && !r.wantEOS
}

func (r *Reader) Read(data []byte) (int, error) {
	block := newReadBlock(data)
	for !block.exhausted() {
		if r.streamEnds() {
			r.streaming = false
			// A explicit END OF STREAM signal needs to be send if
			// some bytes are already read into the buffer. Otherwise
			// the return value is already (0, nil) and we are good to go.
			r.wantEOS = block.readSize > 0
			break
		}
		if r.wantEOS {
			r.wantEOS = false
			return 0, nil
		}
		// Check for EOF in case the reader is closed
		if _, err := r.buffered.Peek(1); err != nil {
			return block.readSize, err
		}
		b, err := r.buffered.ReadByte()
		r.streaming = true
		if err != nil {
			return block.readSize, err
		}
		block.writeByte(b)
		r.traffic++
	}
	return block.readSize, nil
}
