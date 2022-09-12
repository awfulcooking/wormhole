package wormhole

import (
	"errors"
	"sync"
)

type ControllerPool struct {
	controllers map[string]*Controller
	mutex       sync.RWMutex
}

func NewControllerPool() *ControllerPool {
	return &ControllerPool{
		controllers: make(map[string]*Controller),
	}
}

func (p *ControllerPool) Get(name string) *Controller {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.controllers[name]
}

func (p *ControllerPool) SetUniq(name string, c *Controller) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, exists := p.controllers[name]; exists {
		return errors.New("controller name already in use: " + name)
	}
	p.controllers[name] = c
	return nil
}

func (p *ControllerPool) Delete(name string) *Controller {
	p.mutex.Lock()
	defer delete(p.controllers, name)
	defer p.mutex.Unlock()

	return p.controllers[name]
}
