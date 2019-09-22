package common

type SubProc struct {
	ChStop chan struct{}
	ChErr  chan error
}

func NewSubProc() *SubProc {
	return &SubProc{
		ChErr: make(chan error),
	}
}

func (p *SubProc) BaseDispose() {
	close(p.ChErr)
}
