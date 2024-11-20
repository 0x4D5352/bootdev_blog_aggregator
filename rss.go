package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/0x4D5352/bootdev_blog_aggregator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
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

func scrapeFeeds(s *state) error {
	next_feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	feed, err := fetchFeed(context.Background(), next_feed.Url)
	if err != nil {
		return err
	}
	t := time.Now().UTC()
	nt := sql.NullTime{
		Time:  t,
		Valid: true,
	}
	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: nt,
		ID:            next_feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Articles from %s\n", feed.Channel.Title)
	for _, item := range feed.Channel.Item {
		fmt.Printf("Adding %+v\n", item)
		title := sql.NullString{String: item.Title}
		if item.Title != "" {
			title.Valid = true
		}
		description := sql.NullString{String: item.Description}
		if item.Description != "" {
			description.Valid = true
		}
		// fmt.Printf("Checking Pubdate: %s\n", item.PubDate)
		pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			fmt.Println("poof")
			return err
		}
		post, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   t,
			UpdatedAt:   t,
			Title:       title,
			Url:         item.Link,
			Description: description,
			PublishedAt: pubDate,
			FeedID:      next_feed.ID,
		})
		if pqErr, ok := err.(*pq.Error); ok {
			// 23505 == "unique_violation" https://www.postgresql.org/docs/9.3/errcodes-appendix.html
			if pqErr.Code == "23505" {
				fmt.Printf("%s already added!\n", item.Title)
				continue
			}
			return err
		}
		fmt.Printf("Added %+v\n", post)
	}
	return nil
}
