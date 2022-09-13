package wormhole

import (
	"time"
)

type PipeState string

const (
	Pending PipeState = "pending"
	Opened  PipeState = "opened"
	Closed  PipeState = "closed"
)

type PipeDataType int

const (
	DataUTF8 PipeDataType = iota
	DataBinary
)

type PipeClient interface {
	Read() ([]byte, PipeDataType, error)
	Write([]byte, PipeDataType) error
	Close(error) error
}

type Pipe struct {
	ID         PipeID
	Client     PipeClient
	State      PipeState
	resolution chan error
}

func NewPipe(id PipeID, client PipeClient) *Pipe {
	return &Pipe{
		ID:         id,
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

func (p *Pipe) Run(controller *Controller) error {
	if p.State != Opened {
		panic("tried to run non-open pipe")
	}

	for {
		data, dataType, err := p.Client.Read()

		if err != nil {
			return err
		}

		if err := controller.PipeFromClient(p.ID, data, dataType); err != nil {
			return err
		}
	}
}
