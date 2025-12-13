package main

import (
	"context"
	"fmt"
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
)

func main() {
	cfg := configs.MustLoad()

	db := mustOpenDB(cfg)

	compositionRoot := cmd.NewCompositionRoot(cfg, db)
	defer compositionRoot.CloseAll()

	compositionRoot.Logger().Info("bot started")

	startBot(compositionRoot, cfg.TelegramBotToken)
}

func startBot(compositionRoot *cmd.CompositionRoot, token string) {
	bot, err := telegram.NewBot(compositionRoot.Logger(), token, compositionRoot.NewUserRegistrationCommandHandler())
	if err != nil {
		panic(fmt.Sprintf("create bot %s", err))
	}

	bot.HandleUpdates(context.Background())
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
