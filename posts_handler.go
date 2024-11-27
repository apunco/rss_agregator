package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/apunco/go/rss_agregator/internal/database"
)

func getUserPostsHandler(s *state, cmd command, user database.GatorUser) error {
	var limit int
	var err error

	if len(cmd.args) != 1 {
		limit = 2
	} else {
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			fmt.Printf("args %v could not be parse", cmd.args[0])
			return err
		}
	}

	getPostsArgs := database.GetUserPostsParams{
		AddedBy: user.ID,
		Limit:   int32(limit),
	}

	posts, err := s.db.GetUserPosts(context.Background(), getPostsArgs)
	if err != nil {
		fmt.Printf("error getting posts for user %s", user.Name)
	}

	for _, post := range posts {
		fmt.Printf("\nTitle: %s\n", post.Title)
		fmt.Printf("URL: %s\n", post.Url)
		if post.Description.Valid {
			fmt.Printf("Description: %s\n", post.Description.String)
		}
		if post.PublishedAt.Valid {
			fmt.Printf("Published: %s\n", post.PublishedAt.Time.Format("2006-01-02 15:04:05"))
		}
		fmt.Println("-------------------")
	}

	return nil
}
