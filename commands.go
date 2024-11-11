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

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
	fmt.Println("command " + name + " registered for use")
}

func (c *commands) run(s *state, cmd command) error {
	cmdName := cmd.name
	if _, ok := c.commands[cmdName]; !ok {
		return errors.New("command " + cmdName + " is not registered for use")
	}

	err := c.commands[cmd.name](s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing argument <username>")
	}

	username := cmd.args[0]
	err := s.cfg.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Println(username + " set as current user")
	return nil
}
