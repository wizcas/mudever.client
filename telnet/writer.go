package telnet

import "io"

type writer struct {
	dest    io.Writer
	traffic uint64
}

func newWriter(w io.Writer) *writer {
	return &writer{
		dest:    w,
		traffic: 0,
	}
}

func (w *writer) write(data []byte) (int, error) {
	nWritten, err := w.dest.Write(data)
	if err != nil {
		return nWritten, err
	}
	return nWritten, nil
}
