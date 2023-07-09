package bot

import (
	"context"
	"errors"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ttrueno/rl2-final/internal/app/botkit"
	"github.com/ttrueno/rl2-final/internal/app/models"
	"github.com/ttrueno/rl2-final/internal/db/psql"
	"github.com/ttrueno/rl2-final/internal/lib/conv"
)

const (
	nextPageButtonLabel              = "следующая страница."
	prevPageButtonLabel              = "предыдущая страница."
	taskFinishedStatusResponse       = "Готово."
	clientResponseOutOfRangeResponse = "Оу, прости, но ответ должен состоять в множестве допустимых значений((."
	askClientToRespondMessage        = "Вот варианты ответов.Выбери пожалуйста: "
	noRecordsMessage                 = "Оу, прости, но у нас более нет записей с этой страницы, мы перенесём тебя на предыдущую."
)

type compositionsFetcher interface {
	GetAllByComposerID(ctx context.Context, ComposerID string) ([]models.Composition, error)
}

type composersFetcher interface {
	GetAll(ctx context.Context, pageNumber int) ([]models.Composer, error)

	PageLength() int
}

type ViewCmdGetComposers struct {
	bot                 *botkit.Bot
	p                   *paginator
	composersFetcher    composersFetcher
	compositionsFetcher compositionsFetcher
	config              Config
}

type Config struct {
	MaxCompositions int
}

func NewViewCmdComposers(
	bot *botkit.Bot,
	composersFetcher composersFetcher,
	compositionsFetcher compositionsFetcher,
	config Config,
) *ViewCmdGetComposers {
	return &ViewCmdGetComposers{
		bot: bot,
		p: &paginator{
			next:       true,
			prev:       true,
			pageNumber: 1,
			pageLength: composersFetcher.PageLength(),
		},
		composersFetcher:    composersFetcher,
		compositionsFetcher: compositionsFetcher,
		config:              config,
	}
}

func (h *ViewCmdGetComposers) ViewCmdComposers() botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) (err error) {
		h.p.pageNumber = 1
		h.p.reset()

		var (
			msg = tgbotapi.NewMessage(update.FromChat().ID, "")
		)

		defer func() {
			h.p.reset()

			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			msg.Text = taskFinishedStatusResponse
			_, err = bot.Send(msg)
		}()

		for {
			composers, err := h.fetchComposers(ctx, msg, update.FromChat().ID)
			if err != nil {
				return err
			}

			if composers == nil {
				continue
			}

			if err := h.sendKeyboard(ctx, update.FromChat().ID, msg, composers); err != nil {
				return err
			}
			done, err := h.handleIO(ctx, composers)
			if err != nil {
				return err
			}

			if done {
				return nil
			}
		}
	}
}

func (h *ViewCmdGetComposers) handleIO(ctx context.Context, composers []models.Composer) (done bool, err error) {
	update := h.bot.Update(ctx)
	switch update.Message.Text {
	case nextPageButtonLabel:
		h.p.Next()
		return false, nil
	case prevPageButtonLabel:
		h.p.Prev()
		return false, nil
	default:
		break
	}
	msg := tgbotapi.NewMessage(update.FromChat().ID, "")

	num, err := conv.ParseFirstInt64(update.Message.Text)
	if err != nil {
		return false, nil
	}

	if num < 1 || num > int64(len(composers)) {
		msg.Text = clientResponseOutOfRangeResponse
		if err := h.bot.Send(ctx, msg); err != nil {
			return false, err
		}
		return false, nil
	}
	composer := composers[num-1]
	msg.Text = describeComposer(composer)

	if err = h.bot.Send(ctx, msg); err != nil {
		return false, err
	}

	err = h.bot.SendPhotoByURL(ctx, composer.ImageLink, update.FromChat().ID)
	if err != nil {
		return false, err
	}

	var compositions []models.Composition
	err = func() error {
		fetchCompositionsCtx, fetchCompositionsCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer fetchCompositionsCancel()
		compositions, err = h.compositionsFetcher.GetAllByComposerID(fetchCompositionsCtx, composers[num-1].ID)
		if err != nil && !errors.Is(err, psql.ErrNoRecords) {
			return err
		}

		return nil
	}()
	if err != nil {
		return false, err
	}

	compositionsOutput := ""
	compositionNumber := 1
	for _, composition := range compositions {
		compositionsOutput = fmt.Sprintf("%s\n%s", compositionsOutput, describeComposition(compositionNumber, composition))
		compositionNumber++
	}

	if err = h.bot.Send(ctx, tgbotapi.NewMessage(
		update.FromChat().ID,
		compositionsOutput,
	)); err != nil {
		return false, err
	}

	return true, nil
}

func describeComposition(compositionNumber int, composition models.Composition) string {
	res := fmt.Sprintf(`%d: "%s"%s`, compositionNumber, composition.Name, "\n")

	mirrorNumber := 1
	for _, mirror := range composition.Mirrors {
		res = fmt.Sprintf("%s\n\t\t\t%d: %s", res, mirrorNumber, mirror.Link)
		mirrorNumber++
	}
	res = fmt.Sprintf("%s\n", res)

	return res
}

func describeComposer(composer models.Composer) string {
	return fmt.Sprintf("%s %s - %s", composer.FirstName, composer.LastName, composer.Description)
}

func (h *ViewCmdGetComposers) sendKeyboard(ctx context.Context, chatId int64, msg tgbotapi.MessageConfig, composers []models.Composer) error {
	keyboard := h.drawKeyboard(composers)
	msg = tgbotapi.NewMessage(chatId, askClientToRespondMessage)
	msg.ReplyMarkup = keyboard
	if err := h.bot.Send(ctx, msg); err != nil {
		return err
	}

	return nil
}

func (h *ViewCmdGetComposers) fetchComposers(ctx context.Context, msg tgbotapi.MessageConfig, chatID int64) (_ []models.Composer, err error) {
	var composers []models.Composer
	err = func() error {
		fetchCtx, fetchCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer fetchCancel()

		composers, err = h.composersFetcher.GetAll(fetchCtx, h.p.pageNumber)
		if err != nil {
			return err
		}

		return nil
	}()

	if err != nil {
		switch {
		case errors.Is(err, psql.ErrNoRecords):
			if err := h.noRecords(ctx, chatID); err != nil {
				return nil, err
			}

			return nil, nil
		default:
			return nil, err
		}
	}

	h.p.forItems(len(composers))
	return composers, nil
}

func (h *ViewCmdGetComposers) noRecords(ctx context.Context, chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, noRecordsMessage)
	h.p.next = false
	h.p.pageNumber--
	if err := h.bot.Send(ctx, msg); err != nil {
		return err
	}

	return nil
}

func (h *ViewCmdGetComposers) drawKeyboard(composers []models.Composer) *tgbotapi.ReplyKeyboardMarkup {
	keyboard := make([][]tgbotapi.KeyboardButton, 0)

	row := make([]tgbotapi.KeyboardButton, 0)
	for i, composer := range composers {
		buttonLabel := fmt.Sprintf("%d: %s %s", i+1, composer.FirstName, composer.LastName)

		row = append(row, tgbotapi.NewKeyboardButton(buttonLabel))
	}
	keyboard = append(keyboard, row)

	navButtons := make([]tgbotapi.KeyboardButton, 0)
	if h.p.prev {
		navButtons = append(
			navButtons,
			tgbotapi.NewKeyboardButton(prevPageButtonLabel),
		)
	}
	if h.p.next {
		navButtons = append(
			navButtons,
			tgbotapi.NewKeyboardButton(nextPageButtonLabel),
		)
	}

	keyboard = append(keyboard, navButtons)
	return &tgbotapi.ReplyKeyboardMarkup{Keyboard: keyboard, ResizeKeyboard: true}
}

type paginator struct {
	prev       bool
	next       bool
	pageLength int
	pageNumber int
}

func (p *paginator) Next() {
	p.pageNumber++
	p.prev = true
}

func (p *paginator) Prev() {
	p.pageNumber--
	p.next = true
}

func (p *paginator) reset() {
	p.next = true
	p.prev = true
}

func (p *paginator) forItems(itemsNum int) {
	if itemsNum < p.pageLength {
		p.next = false
	}

	if p.pageNumber == 1 {
		p.prev = false
	}
}
