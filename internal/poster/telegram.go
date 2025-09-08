package poster

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/FANIMAN/cncg-bot/internal/store"
	"github.com/FANIMAN/cncg-bot/internal/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TG struct {
	Bot       *tgbotapi.BotAPI
	ChannelID string
	ParseMode string // "MarkdownV2"
}

func New(botToken, channelID, parseMode string) (*TG, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}
	return &TG{Bot: bot, ChannelID: channelID, ParseMode: parseMode}, nil
}

func (t *TG) PostItem(ctx context.Context, it store.Item) error {
	title := util.EscapeTelegram(it.Title)
	url := util.EscapeTelegram(it.URL)
	source := util.EscapeTelegram(it.Source)
	sum := util.EscapeTelegram(it.Summary)

	// Title as a clickable link, then 2–3 sentence summary + source attribution + tags
	text := fmt.Sprintf("[*%s*](%s)\n\n%s\n\n_Source:_ %s", title, url, sum, source)

	channelId, err := strconv.Atoi(os.Getenv("CHANNEL_ID"))
	if err != nil {
		return err
	}
	msg := tgbotapi.MessageConfig{
		BaseChat:              tgbotapi.BaseChat{ChatID: int64(channelId)},
		Text:                  text,
		ParseMode:             t.ParseMode,
		DisableWebPagePreview: false,
	}
	_, err = t.Bot.Send(msg)
	return err
}
