package wormhole

import "errors"

type NameGenerator func() string

func (g NameGenerator) Assign(pool *ControllerPool, c *Controller, attempts int) error {
	for i := 0; i < attempts; i++ {
		name := g()
		if pool.SetUniq(name, c) {
			c.Name = name
			return nil
		}
	}

	return ErrNameGeneratorExhausted
}

var ErrNameGeneratorExhausted = errors.New("name generator exhaused. no unique name after max attempts")

func StaticNameGenerator(name string) NameGenerator {
	return func() string {
		return name
	}
}
