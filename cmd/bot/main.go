package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/FANIMAN/cncg-bot/internal/config"
	"github.com/FANIMAN/cncg-bot/internal/core"
	"github.com/FANIMAN/cncg-bot/internal/poster"
	"github.com/FANIMAN/cncg-bot/internal/store"
)

func handler(ctx context.Context) error {

	cfg := config.Load()

	db, err := store.Open(cfg.DBPath)
	if err != nil {
		log.Fatal("db:", err)
	}

	p := core.Pipeline{
		Filters: core.Filters{
			MaxAgeDays: cfg.Filters.MaxAgeDays,
			MinScore:   cfg.Filters.MinScore,
			Positive:   cfg.Keywords.Positive,
			Negative:   cfg.Keywords.Negative,
		},
		DB: db,
	}
	for _, s := range cfg.Sources {
		p.Sources = append(p.Sources, core.SourceCfg{
			Name: s.Name, Type: s.Type, URL: s.URL, Weight: s.Weight, Tags: s.Tags,
		})
	}

	tg, err := poster.New(cfg.Telegram.BotToken, cfg.Telegram.ChannelID, cfg.Telegram.ParseMode)
	if err != nil {
		return err
	}

	// Run pipeline once
	err = p.RunOnce(ctx)
	if err != nil {
		db.LogError(ctx, "pipeline", err.Error())
		return err
	}

	// Send next batch
	items, err := db.NextUnposted(ctx, cfg.Filters.MinScore, cfg.Scheduler.BatchSize)
	if err != nil {
		db.LogError(ctx, "schedule:select", err.Error())
		return err
	}

	for _, it := range items {
		if err := tg.PostItem(ctx, it); err != nil {
			db.LogError(ctx, "telegram:send", err.Error())
			continue
		}
		if err := db.MarkPosted(ctx, it.ID); err != nil {
			log.Println("mark posted error:", err)
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
