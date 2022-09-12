package wormhole

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"sync/atomic"
)

type Controller struct {
	context context.Context
	Host    ControllerHost

	pipes map[PipeID]*Pipe
	mutex sync.RWMutex

	nextPipeID uint64
}

func NewController(ctx context.Context, host ControllerHost) *Controller {
	return &Controller{
		context: ctx,
		Host:    host,
	}
}

func (c *Controller) ProcessNext() error {
	packet, err := c.Host.ReadControllerPacket()
	if err != nil {
		return err
	}
	return c.Process(packet)
}

var ErrAlienPacket = errors.New("packet format not recognised")
var ErrPipeNotFound = errors.New("pipe not found")

func (c *Controller) Process(packet ControllerPacket) error {
	if packet.PipeID != nil {
		pipe, ok := c.pipes[*packet.PipeID]
		if !ok {
			return ErrPipeNotFound
		}

		if packet.State != nil {
			switch *packet.State {
			case "opened":
				pipe.Fulfil()
			case "closed":
				pipe.Close()
			case "error":
				pipe.Reject()
			default:
				return ErrAlienPacket
			}

			return nil
		} else if packet.Message != nil {
			pipe.Client.Write(packet.Message)
		}
	}

	return ErrAlienPacket
}

func (c *Controller) RequestPipe(ctx context.Context, client io.ReadWriteCloser) (*Pipe, error) {
	pipe := NewPipe(client)
	id := c.NextPipeID()

	c.mutex.Lock()
	c.pipes[id] = pipe
	c.mutex.Unlock()

	err := pipe.Wait()
	return pipe, err
}

func (c *Controller) NextPipeID() PipeID {
	return PipeID(atomic.AddUint64(&c.nextPipeID, 1))
}

func (c *Controller) AcceptClient(client net.Conn) error {
	c.mutex.Lock()
	c.pipes = append(c.pipes, &Pipe{})
	c.mutex.Unlock()
	return nil
}
