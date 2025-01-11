package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"

	"github.com/mashfeii/gator/internal/domain"
)

func FetchFeed(ctx context.Context, url string) (*domain.RSSFeed, error) {
	client := http.DefaultClient

	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("unable to create http request: %w", err)
	}

	req.Header.Set("User-Agent", "gator")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to make http request: %w", err)
	}
	defer resp.Body.Close()

	var feed domain.RSSFeed

	err = xml.NewDecoder(resp.Body).Decode(&feed)
	if err != nil {
		return nil, fmt.Errorf("unable to decode xml: %w", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	return &feed, nil
}
