package main

import (
	"errors"
)

type command struct{
	name string
	Args []string
}

type commands struct{
	Handler  map[string]func(*state, command) error
}

func (c *commands) run (s *state, cmd command) error {
		handler, exists := c.Handler[cmd.name]
		if exists != true {
		return errors.New("no command like this present in commands")
	}
		return handler(s, cmd) 
}

func (c *commands) register(name string, f func(*state, command) error){
	c.Handler[name] = f
}
