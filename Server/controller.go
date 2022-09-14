package wormhole

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type Controller struct {
	context context.Context
	Host    ControllerHost
	Name    string

	pipes map[PipeID]*Pipe
	mutex sync.RWMutex

	nextPipeID uint64
}

func NewController(ctx context.Context, host ControllerHost) *Controller {
	return &Controller{
		context: ctx,
		Host:    host,
		pipes:   make(map[PipeID]*Pipe),
	}
}

func (c *Controller) SendMeta(meta ControllerMeta) error {
	return c.Host.WriteControllerPacket(ControllerPacket{
		Meta: &meta,
	})
}

func (c *Controller) SendWelcome() error {
	return c.SendMeta(ControllerMeta{
		Name: c.Name,
	})
}

func (c *Controller) ProcessNext() error {
	println("asked to read")
	packet, err := c.Host.ReadControllerPacket()
	if err != nil {
		return err
	}
	println("successful controller packet read")
	return c.Process(packet)
}

var ErrAlienPacket = errors.New("packet format not recognised")
var ErrPipeNotFound = errors.New("pipe not found")

func (c *Controller) Process(packet ControllerPacket) error {
	pipe, ok := c.pipes[packet.PipeID]
	if !ok {
		return ErrPipeNotFound
	}

	if packet.State != "" {
		switch packet.State {
		case Opened:
			pipe.Fulfil()
		case Closed:
			pipe.Close()
		case "error":
			msg := "host reported pipe error"
			if packet.Data != nil {
				msg = string(*packet.Data)
			}
			pipe.Reject(errors.New(msg))
		default:
			return ErrAlienPacket
		}

		return nil
	}

	if packet.Data != nil {
		return pipe.Client.Write([]byte(*packet.Data), DataUTF8)
	}

	return ErrAlienPacket
}

func (c *Controller) RequestPipe(ctx context.Context, client PipeClient, subprotocol string) (*Pipe, error) {
	pipe := NewPipe(c.NextPipeID(), client)

	c.mutex.Lock()
	c.pipes[pipe.ID] = pipe
	c.mutex.Unlock()

	var state = Pending

	c.Host.WriteControllerPacket(ControllerPacket{
		PipeID:      pipe.ID,
		State:       state,
		Subprotocol: subprotocol,
	})

	err := pipe.Wait(10 * time.Second)
	return pipe, err
}

func (c *Controller) NextPipeID() PipeID {
	return PipeID(atomic.AddUint64(&c.nextPipeID, 1))
}

func (c *Controller) PipeFromClient(pipeID PipeID, data []byte, dataType DataType) error {
	dataStr := string(data)
	return c.Host.WriteControllerPacket(ControllerPacket{
		PipeID:   pipeID,
		Data:     &dataStr,
		DataType: dataType,
	})
}

func (c *Controller) Shutdown() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, pipe := range c.pipes {
		pipe.Close()
	}

	return c.Host.Close()
}
