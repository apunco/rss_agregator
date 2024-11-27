package main

import (
	"context"
	"errors"
	"log"
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
	//fmt.Println("command " + name + " registered for use")
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

func resetHandler(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		log.Printf("error deleting users %s", err)
		return err
	}

	return nil
}
