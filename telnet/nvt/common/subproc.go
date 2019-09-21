package common

type SubProc struct {
	ChStop chan struct{}
	ChErr  chan error
}

func NewSubProc() *SubProc {
	return &SubProc{
		ChStop: make(chan struct{}),
		ChErr:  make(chan error),
	}
}

func (p *SubProc) Stop() {
	p.ChStop <- struct{}{}
	close(p.ChStop)
}

func (p *SubProc) BaseDispose() {
	close(p.ChErr)
}
