package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"log"
	"net/http"
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
