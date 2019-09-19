package stream

import "io"

// Writer writes bytes into a writer that keep alive
// (i.e. not necessarily throws EOF when there's no more data)
type Writer struct {
	dest    io.Writer
	traffic uint64
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		dest:    w,
		traffic: 0,
	}
}

func (w *Writer) Write(data []byte) (int, error) {
	nWritten, err := w.dest.Write(data)
	if err != nil {
		return nWritten, err
	}
	return nWritten, nil
}
