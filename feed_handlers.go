package main

import (
	"context"
	"fmt"
)

func getFeedHandler(s *state, cmd command) error {

	feed, err := fetchFeed(context.Background(), "https://wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", feed)
	return nil
}

func addFeedHandler(s *state, cmd command) error {

	return nil
}
