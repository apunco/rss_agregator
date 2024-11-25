package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/apunco/go/rss_agregator/internal/database"
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
	if len(cmd.args) != 2 {
		return errors.New("missing argument <name> <url>")
	}

	user, err := s.db.GetUserByName(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		fmt.Printf("error getting current user %s", err)
		return err
	}

	args := database.AddFeedParams{
		Name:    cmd.args[0],
		Url:     cmd.args[1],
		AddedBy: user.ID,
	}

	feed, err := s.db.AddFeed(context.Background(), args)
	if err != nil {
		fmt.Printf("error adding feed %s", err)
		return err
	}

	feedFollowArgs := database.AddFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	if _, err := s.db.AddFeedFollow(context.Background(), feedFollowArgs); err != nil {
		fmt.Printf("error creating feed follow %s", err)
		return err
	}

	fmt.Println(feed)
	return nil
}

func getFeedsHandler(s *state, cmd command) error {
	type printableFeed struct {
		userName string
		name     string
		url      string
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		fmt.Printf("error getting feeds %s", err)
	}

	for _, feed := range feeds {
		userName, err := s.db.GetUserName(context.Background(), feed.AddedBy)
		if err != nil {
			fmt.Printf("error getting user name %s", err)
		}

		printFeed := printableFeed{
			userName: userName,
			name:     feed.Name,
			url:      feed.Url,
		}

		fmt.Println(printFeed)
	}

	return nil
}

func addFeedFollowHandler(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("missing argument <name> <url>")
	}

	user, err := s.db.GetUserByName(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		fmt.Printf("error getting current user %s", err)
		return err
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		fmt.Printf("error getting feed with url %s", err)
		return err
	}

	args := database.AddFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	feedFollowRow, err := s.db.AddFeedFollow(context.Background(), args)
	if err != nil {
		fmt.Printf("error creating feed follow %s", err)
		return err
	}

	fmt.Println(feedFollowRow)

	return nil
}

func getFollowingForUserHandler(s *state, cmd command) error {

	userName := s.cfg.CurrentUserName
	user, err := s.db.GetUserByName(context.Background(), userName)
	if err != nil {
		fmt.Printf("error getting user id %s", err)
		return err
	}

	feeds, err := s.db.GetFeedsForUser(context.Background(), user.ID)
	if err != nil {
		fmt.Printf("error getting feeds for user %s", err)
		return err
	}

	for _, feed := range feeds {
		fmt.Println(feed)
	}

	return nil
}
