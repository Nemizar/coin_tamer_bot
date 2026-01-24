package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Nemizar/coin_tamer_bot/internal/adapters/in/telegram"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Nemizar/coin_tamer_bot/cmd"
	"github.com/Nemizar/coin_tamer_bot/configs"
)

const (
	setMaxIdleConns    = 2
	setMaxOpenConns    = 5
	setConnMaxLifetime = 10 * time.Minute
	setConnMaxIdleTime = 10 * time.Minute
	shutdownTimeout    = 30 * time.Second
)

func main() {
	cfg := configs.MustLoad()

	db := mustOpenDB(cfg)
	compositionRoot := cmd.NewCompositionRoot(cfg, db)
	defer compositionRoot.CloseAll()

	logger := compositionRoot.Logger()

	logger.Info("bot starting", "env", cfg.ENV)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	errCh := make(chan error, 1)
	done := make(chan struct{})

	go func() {
		defer close(done)

		if err := startBot(ctx, compositionRoot, cfg); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case err := <-errCh:
		logger.Error("bot stopped with error", "err", err)
	}

	shutdownTimer := time.NewTimer(shutdownTimeout)
	defer shutdownTimer.Stop()

	select {
	case <-done:
		logger.Info("bot stopped gracefully")
	case <-shutdownTimer.C:
		logger.Error("shutdown timeout exceeded")
	}
}

func startBot(
	ctx context.Context,
	compositionRoot *cmd.CompositionRoot,
	cfg configs.Config,
) error {
	bot, err := telegram.NewBot(
		compositionRoot.Logger(),
		cfg.TelegramBotToken,
		cfg.AllowedChatIDs,
		compositionRoot.NewUserRegistrationCommandHandler(),
		compositionRoot.NewCreateDefaultCategoryCommandHandler(),
		compositionRoot.NewCreateTransactionCommandHandler(),
		compositionRoot.NewGetCategoriesByTypeQueryHandler(),
		compositionRoot.NewGetUserQueryHandler(),
	)
	if err != nil {
		return fmt.Errorf("create bot: %w", err)
	}

	bot.HandleUpdates(ctx)

	return nil
}

func mustOpenDB(cfg configs.Config) *sqlx.DB {
	var (
		db  *sqlx.DB
		err error
	)

	db, err = sqlx.Open("pgx", cfg.DBDSNString())
	if err != nil {
		panic(fmt.Sprintf("build db client: %s", err))
	}

	db.SetMaxIdleConns(setMaxIdleConns)
	db.SetMaxOpenConns(setMaxOpenConns)
	db.SetConnMaxLifetime(setConnMaxLifetime)
	db.SetConnMaxIdleTime(setConnMaxIdleTime)

	return db
}
