package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ttrueno/rl2-final/internal/app/botkit"
)

const HelpResponse = `Привет.

Меня зовут Саня.
Я распространяю классическую музыку на своём районе.
К сожалению на данный момент ты можешь получить информацию только о том, где и как послушать того или иного композитора.

Вот доступные опции:
/start - Чтобы начать разговор
/help - Чтобы получить инструкцию
/composers - Чтобы получить возможность получить источники на произведения того или иного композитора, йееей.
`

func ViewCmdHelp() botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, HelpResponse)); err != nil {
			return err
		}

		return nil
	}
}
