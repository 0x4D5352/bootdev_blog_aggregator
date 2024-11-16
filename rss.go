package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"strings"
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
	client := &http.Client{}
	reader := strings.NewReader("")
	request, err := http.NewRequestWithContext(ctx, "GET", feedURL, reader)
	if err != nil {
		return &RSSFeed{}, err
	}
	request.Header.Add("User-Agent", "gator")
	response, err := client.Do(request)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return &RSSFeed{}, err
	}
	var feed RSSFeed
	xml.Unmarshal(body, &feed)
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}
	return &feed, nil
}
