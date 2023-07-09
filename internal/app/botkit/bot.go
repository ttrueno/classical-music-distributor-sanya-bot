package botkit

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gookit/slog"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	cmdViews map[string]ViewFunc
	updates  tgbotapi.UpdatesChannel
}

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error

func New(api *tgbotapi.BotAPI) *Bot {
	return &Bot{
		api: api,
	}
}

func (b *Bot) RegisterCmdView(cmd string, view ViewFunc) {
	if b.cmdViews == nil {
		b.cmdViews = make(map[string]ViewFunc)
	}

	b.cmdViews[cmd] = view
}

func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	b.updates = b.api.GetUpdatesChan(u)

	for {
		select {
		case update := <-b.updates:
			updateCtx, updateCancel := context.WithTimeout(context.Background(), 5*time.Second)
			b.handleViewUpdate(updateCtx, update)
			updateCancel()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *Bot) Update(ctx context.Context) tgbotapi.Update {
	return <-b.updates
}

func (b *Bot) Send(ctx context.Context, msg tgbotapi.MessageConfig) error {
	_, err := b.api.Send(msg)
	return err
}

func (b *Bot) SendPhotoByURL(ctx context.Context, url string, chatID int64) error {
	if url == "" {
		return nil
	}

	img := tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL(url))

	mediaGroup := tgbotapi.NewMediaGroup(
		chatID,
		[]interface{}{
			img,
		},
	)
	if _, err := b.api.SendMediaGroup(mediaGroup); err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleViewUpdate(ctx context.Context, update tgbotapi.Update) {
	defer func() {
		if p := recover(); p != nil {
			slog.Panicf("panic recovered: %v", p)
		}
	}()

	var view ViewFunc

	if !update.Message.IsCommand() {
		return
	}

	cmd := update.Message.Command()

	cmdView, ok := b.cmdViews[cmd]
	if !ok {
		return
	}

	view = cmdView
	if err := view(ctx, b.api, update); err != nil {
		slog.Errorf("failed to handle update: %v", err)

		b.errorMsg(ctx, update)
	}
}

func (b *Bot) errorMsg(ctx context.Context, update tgbotapi.Update) {
	if _, err := b.api.Send(
		tgbotapi.NewMessage(update.Message.Chat.ID, "internal error"),
	); err != nil {
		slog.Errorf("failed to send error message: %v", err)
	}
}
