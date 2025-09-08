package fetch

import (
	"context"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

type RSSSource struct {
	Name   string
	URL    string
	Tags   []string
	Weight float64
}

type Item struct {
	Source      string
	Title       string
	URL         string
	PublishedAt time.Time
	Summary     string
	Tags        []string
}

func FetchRSS(ctx context.Context, src RSSSource) ([]Item, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(src.URL, ctx)
	if err != nil { return nil, err }

	var items []Item
	for _, e := range feed.Items {
		u := strings.TrimSpace(firstNonEmpty(e.Link, e.GUID))
		t := strings.TrimSpace(e.Title)
		pub := time.Now().UTC()
		if e.PublishedParsed != nil { pub = e.PublishedParsed.UTC() }
		if e.UpdatedParsed != nil { pub = e.UpdatedParsed.UTC() }
		summary := strings.TrimSpace(firstNonEmpty(e.Description, e.Content))
		items = append(items, Item{
			Source: src.Name, Title: t, URL: u, PublishedAt: pub, Summary: summary, Tags: src.Tags,
		})
	}
	return items, nil
}

func firstNonEmpty(ss ...string) string {
	for _, s := range ss {
		if strings.TrimSpace(s) != "" { return s }
	}
	return ""
}
