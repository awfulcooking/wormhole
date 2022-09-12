package wormhole

import (
	"io"
	"time"
)

type PipeState int

const (
	Pending PipeState = iota
	Opened
	Closed
)

type Pipe struct {
	Client     io.ReadWriteCloser
	State      PipeState
	resolution chan error
}

func NewPipe(client io.ReadWriteCloser) *Pipe {
	return &Pipe{
		Client:     client,
		State:      Pending,
		resolution: make(chan error, 1),
	}
}

func (p *Pipe) Wait(deadline time.Duration) error {
	select {
	case err := <-p.resolution:
		return err
	case <-time.After(deadline):
		return <-p.resolution
	}
}

func (p *Pipe) resolve(err error) {
	p.resolution <- err
	close(p.resolution)
}

func (p *Pipe) Fulfil() {
	p.State = Opened
	p.resolve(nil)
}

func (p *Pipe) Reject(err error) {
	p.State = Closed
	p.resolve(err)
}

func (p *Pipe) Close() {
	p.State = Closed
}
