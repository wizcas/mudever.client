package common

type SubProc struct {
	chErr chan error
}

func NewSubProc() *SubProc {
	return &SubProc{
		chErr: make(chan error),
	}
}

func (p *SubProc) GotError(err error) {
	if p.chErr == nil {
		p.chErr = make(chan error)
	}
	p.chErr <- err
}

func (p *SubProc) Err() <-chan error {
	return p.chErr
}

func (p *SubProc) BaseDispose() {
	close(p.chErr)
}
