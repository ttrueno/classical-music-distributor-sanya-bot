package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ttrueno/rl2-final/internal/app/botkit"
)

const StartResponse string = `Привет.

Меня зовут Саня.
Я распространяю классическую музыку на своём районе.
Отправь команду /help, чтобы узнать то, на что я способен.
`

func ViewCmdStart() botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, StartResponse)); err != nil {
			return err
		}

		return nil
	}
}
