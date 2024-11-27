package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/apunco/go/rss_agregator/internal/database"
)

func getFeedHandler(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("missing argument <time_between_reqs>")
	}

	fmt.Printf("Collecting feeds every %s", cmd.args[0])
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		fmt.Printf("error parsing time between requests %s", err)
		return err
	}

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		if err = scrapeFeeds(s); err != nil {
			fmt.Printf("error scraping feed %s", err)
			return err
		}
	}
}

func addFeedHandler(s *state, cmd command, user database.GatorUser) error {
	if len(cmd.args) != 2 {
		return errors.New("missing argument <name> <url>")
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

func addFeedFollowHandler(s *state, cmd command, user database.GatorUser) error {
	if len(cmd.args) != 1 {
		return errors.New("missing argument <name> <url>")
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

func getFollowingForUserHandler(s *state, cmd command, user database.GatorUser) error {
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

func unfollowFeedHandler(s *state, cmd command, user database.GatorUser) error {
	if len(cmd.args) != 1 {
		return errors.New("missing argument <feed url>")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		fmt.Printf("error getting feed follow %s", err)
		return err
	}

	args := database.DeleteFeedFollowForUserParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	if err = s.db.DeleteFeedFollowForUser(context.Background(), args); err != nil {
		fmt.Printf("error deleting feed follow %s", err)
		return err
	}

	return nil
}
