package core

import (
	"context"
	"strings"
	"time"

	"github.com/FANIMAN/cncg-bot/internal/fetch"
	"github.com/FANIMAN/cncg-bot/internal/store"
	"github.com/FANIMAN/cncg-bot/internal/util"
)

type Filters struct {
	MaxAgeDays int
	MinScore   float64
	Positive   []string
	Negative   []string
}

type SourceCfg struct {
	Name   string
	Type   string
	URL    string
	Weight float64
	Tags   []string
}

type Pipeline struct {
	Filters Filters
	Sources []SourceCfg
	DB      *store.Store
}

func (p *Pipeline) RunOnce(ctx context.Context) error {
	now := time.Now().UTC()
	cutoff := now.AddDate(0, 0, -p.Filters.MaxAgeDays)

	for _, src := range p.Sources {
		switch src.Type {
		case "rss":
			items, err := fetch.FetchRSS(ctx, fetch.RSSSource{
				Name: src.Name, URL: src.URL, Tags: src.Tags, Weight: src.Weight,
			})
			if err != nil {
				p.DB.LogError(ctx, "fetch:rss", src.Name+" : "+err.Error())
				continue
			}

			for _, it := range items {
				if it.PublishedAt.Before(cutoff) {
					continue
				}
				url := util.CanonURL(it.URL)
				title := strings.TrimSpace(it.Title)
				if title == "" || url == "" {
					continue
				}

				rawSum := Summarize(it.Summary, 3)
				score := p.scoreItem(title+" "+rawSum, src.Weight)
				rec := store.Item{
					Source: src.Name, Title: title, URL: url,
					Summary:     rawSum,
					PublishedAt: it.PublishedAt,
					Tags:        strings.Join(src.Tags, ","),
					Hash:        store.Hash(url, title),
					Score:       score,
				}
				if _, err := p.DB.InsertIfNew(ctx, rec); err != nil {
					p.DB.LogError(ctx, "db:insert", err.Error())
				}
			}
		default:
			p.DB.LogError(ctx, "fetch", "unsupported source type: "+src.Type)
		}
	}
	return nil
}

func (p *Pipeline) scoreItem(text string, sourceWeight float64) float64 {
	t := strings.ToLower(text)
	score := 0.2 * sourceWeight
	for _, kw := range p.Filters.Positive {
		if strings.Contains(t, strings.ToLower(kw)) {
			score += 0.1
		}
	}
	for _, kw := range p.Filters.Negative {
		if strings.Contains(t, strings.ToLower(kw)) {
			score -= 0.2
		}
	}
	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}
	if score >= 0.2 {
		return 1
	}

	return score
}
