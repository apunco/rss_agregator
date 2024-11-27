package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/apunco/go/rss_agregator/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {

	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		log.Printf("error creating req context %s", err)
		return nil, err
	}

	req.Header.Add("User-Agent", "rss_gator")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error while calling req %s", err)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error getting resp body %s", err)
		return nil, err
	}

	rssFeed := RSSFeed{}
	if err = xml.Unmarshal(body, &rssFeed); err != nil {
		log.Printf("error unmarshaling rss feed %s", err)
		return nil, err
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)

	return &rssFeed, nil
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Printf("error getting next fetchable feed %s", err)
		return err
	}

	if err = s.db.MarkedFeedFetched(context.Background(), feed.ID); err != nil {
		fmt.Printf("error marking feed as fetched %s", err)
		return err
	}

	fetchedFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	var createPostArgs database.CreatePostParams

	for _, item := range fetchedFeed.Channel.Item {
		var insertTime sql.NullTime
		parsedTime, err := parseTime(item.PubDate)
		if err != nil {
			insertTime = sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			}
		} else {
			insertTime = sql.NullTime{
				Time:  parsedTime,
				Valid: true,
			}
		}

		if utf8.RuneCountInString(item.Description) > 200 {
			item.Description = string([]rune(item.Description)[:200])
		}

		createPostArgs = database.CreatePostParams{
			Title: item.Title,
			Url:   item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != ""},
			PublishedAt: insertTime,
			FeedID:      feed.ID,
		}

		if err = s.db.CreatePost(context.Background(), createPostArgs); err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}

			log.Printf("error creating post: %v", err)
		}

		log.Printf("added post %s to DB", item.Title)
	}

	return nil
}

func parseTime(timeString string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		time.RFC1123,
		time.RFC822,
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"Monday, January 02, 2006 15:04:05 MST",
	}

	var parsedTime time.Time
	var err error

	for _, format := range formats {
		parsedTime, err = time.Parse(format, timeString)
		if err == nil {
			break
		}
	}

	if parsedTime.IsZero() {
		return time.Time{}, errors.New("unable to parse time string " + timeString)
	}

	return parsedTime.UTC(), nil
}
