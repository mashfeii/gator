package domain

import (
	"fmt"
)

type Handler func(*State, *Command) error

type Commands struct {
	Handlers map[string]Handler
}

func (c *Commands) Register(name string, handler Handler) {
	c.Handlers[name] = handler
}

func (c *Commands) Run(s *State, cmd *Command) error {
	command, ok := c.Handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}

	return command(s, cmd)
}
