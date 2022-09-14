package wormhole

import (
	"sync"
)

type ControllerPool struct {
	controllers map[string]*Controller
	mutex       sync.RWMutex

	NameGenerator func() string
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

func (p *ControllerPool) SetUniq(name string, c *Controller) (success bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, found := p.controllers[name]; found {
		return false
	}
	p.controllers[name] = c
	return true
}

func (p *ControllerPool) Delete(name string) *Controller {
	p.mutex.Lock()
	defer delete(p.controllers, name)
	defer p.mutex.Unlock()

	return p.controllers[name]
}
