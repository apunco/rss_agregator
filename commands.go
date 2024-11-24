package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/apunco/go/rss_agregator/internal/database"
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

	_, err := s.db.GetUserByName(context.Background(), username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("user %s doesn't exists", username)
			return err
		}
		log.Printf("something went wrong %s", err)
		return err
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Println(username + " set as current user")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing argument <username>")
	}

	params := database.CreateUserParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}

	userName := cmd.args[0]

	user, err := s.db.GetUserByName(context.Background(), userName)
	if err == nil {
		log.Printf("user already exists")
		return errors.New("user already exists")
	} else if !errors.Is(err, sql.ErrNoRows) {
		log.Printf("something went wrong %s", err)
		return err
	}

	user, err = s.db.CreateUser(context.Background(), params)
	if err != nil {
		log.Printf("error creating user %s", err)
		return err
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		log.Printf("error setting user %s", err)
		return err
	}

	log.Printf("user %s created", user.Name)

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

func getUsersHandler(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		log.Printf("error getting users %s", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func getFeedHandler(s *state, cmd command) error {

	feed, err := fetchFeed(context.Background(), "https://wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", feed)
	return nil
}
