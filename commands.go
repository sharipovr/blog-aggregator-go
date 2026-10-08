package main

import "errors"

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if f, ok := c.handlers[cmd.name]; ok {
		return f(s, cmd)
	}
	return errors.New("no handler found registered for this name")
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}
