// Package main предоставляет точку входа для приложения бота Coin Tamer.
// Выполняет инициализацию бота, настройку зависимостей и запуск обработки сообщений.
package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Nemizar/coin_tamer_bot/configs"
	"github.com/Nemizar/coin_tamer_bot/internal/adapters/out/sl"
	"github.com/Nemizar/coin_tamer_bot/internal/adapters/out/sl/handlers/slogpretty"
	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"
)

func main() {
	cfg := configs.MustLoad()

	logger := setupLogger(cfg)

	logger.Info(fmt.Sprintf("DB: %s@%s:%s\n", cfg.DBUser, cfg.DBHost, cfg.DBPort))
}

func setupLogger(c configs.Config) ports.Logger {
	var log ports.Logger

	if c.IsProd() {
		log = sl.NewSlogLogger(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
	} else {
		log = setupPrettySlog()
	}

	return log
}

func setupPrettySlog() ports.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return sl.NewSlogLogger(slog.New(handler))
}
