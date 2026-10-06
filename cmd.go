package main

import (
	"errors"
	"fmt"
)

type command struct {
	name string
	args []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("Usage: login <username>")
	}
	s.config.SetUser(cmd.args[0])
	fmt.Printf("User: '%s' has been set\n", cmd.args[0])
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	if fn, ok := c.commands[cmd.name]; ok {
		err := fn(s, cmd)
		if err != nil {
			return err
		}
		return nil
	} else {
		return fmt.Errorf("Command not found: %s", cmd.name)
	}
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}
