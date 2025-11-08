package main

import (
	"fmt"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	cmdMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	input, ok := c.cmdMap[cmd.Name]
	if !ok {
		return fmt.Errorf("command %s does not exist", cmd.Name)
	}
	return input(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	if c.cmdMap == nil {
		c.cmdMap = make(map[string]func(*state, command) error)
	}
	if _, ok := c.cmdMap[name]; ok {
		fmt.Printf("command %s already exists\n", name)
		return
	}
	c.cmdMap[name] = f
}
