package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", "gator")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching response: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("received bad status code: %v", resp.StatusCode)
	}

	xmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var rssFeed RSSFeed
	err = xml.Unmarshal(xmlData, &rssFeed)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling xml data: %w", err)
	}

	ch := &rssFeed.Channel
	ch.Title = html.UnescapeString(ch.Title)
	ch.Description = html.UnescapeString(ch.Description)
	for i := range ch.Item {
		ch.Item[i].Title = html.UnescapeString(ch.Item[i].Title)
		ch.Item[i].Description = html.UnescapeString(ch.Item[i].Description)
	}
	return &rssFeed, nil
}
