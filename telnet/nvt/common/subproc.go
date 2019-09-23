package common

import "context"

type SubProc interface {
	Run(ctx context.Context)
	Err() <-chan error
	Stopped() <-chan struct{}
}

type BaseSubProc struct {
	chErr     chan error
	chStopped chan struct{}
}

func NewBaseSubProc() *BaseSubProc {
	return &BaseSubProc{
		chErr:     make(chan error),
		chStopped: make(chan struct{}),
	}
}

func (p *BaseSubProc) GotError(err error) {
	if p.chErr == nil {
		p.chErr = make(chan error)
	}
	p.chErr <- err
}

func (p *BaseSubProc) Err() <-chan error {
	return p.chErr
}

func (p *BaseSubProc) Stopped() <-chan struct{} {
	return p.chStopped
}

func (p *BaseSubProc) BaseDispose() {
	close(p.chErr)
	close(p.chStopped)
}
