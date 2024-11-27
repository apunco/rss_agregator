package main

import (
	"context"

	database "github.com/apunco/go/rss_agregator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.GatorUser) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUserByName(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
