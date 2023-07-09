package main

import (
	"context"
	"errors"
	"os/signal"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gookit/slog"
	"github.com/ttrueno/rl2-final/internal/app/bot"
	"github.com/ttrueno/rl2-final/internal/app/botkit"
	"github.com/ttrueno/rl2-final/internal/app/service/composers"
	"github.com/ttrueno/rl2-final/internal/app/service/compositions"
	composersStorage "github.com/ttrueno/rl2-final/internal/app/storage/composers/psql"
	compositionsStorage "github.com/ttrueno/rl2-final/internal/app/storage/compositions/psql"
	psqlClient "github.com/ttrueno/rl2-final/internal/db/psql"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		slog.Fatalf("failed to load config: %v", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := psqlClient.Connect(ctx, cfg.DbConnConfig)
	if err != nil {
		slog.Fatalf("failed to connect to database: %v", err)
		return
	}
	defer db.Close(ctx)

	composerStorage := composersStorage.New(db.Conn)
	compositionsStorage := compositionsStorage.New(db.Conn)

	composersService := composers.NewService(composerStorage, *composers.NewConfig(3))
	compositionsService := compositions.New(compositionsStorage, *compositions.NewConfig(5))

	botAPI, err := tgbotapi.NewBotAPI(cfg.BotConfig.BotApiToken)
	if err != nil {
		slog.Fatalf("failed to create the bot: %v", err)
		return
	}

	newsBot := botkit.New(botAPI)
	newsBot.RegisterCmdView("start", bot.ViewCmdStart())
	newsBot.RegisterCmdView("help", bot.ViewCmdHelp())
	newsBot.RegisterCmdView("composers",
		bot.NewViewCmdComposers(
			newsBot,
			composersService,
			compositionsService,
			bot.Config{
				MaxCompositions: 10,
			},
		).ViewCmdComposers(),
	)

	botCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL)
	defer stop()

	slog.Println("bot started working: T.T")
	if err := newsBot.Run(botCtx); err != nil {
		if !errors.Is(err, context.Canceled) {
			slog.Errorf("failed to start bot: %v", err)
			return
		}

		slog.Println("bot stopped")
	}
}
