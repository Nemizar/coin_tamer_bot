package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type commandHandler interface {
	command() string
	handle(ctx context.Context, update tgbotapi.Update) error
}

type messageHandler interface {
	canHandle(update tgbotapi.Update) bool
	handle(ctx context.Context, update tgbotapi.Update) error
}

type callbackHandler interface {
	callbackPrefix() string
	handle(ctx context.Context, update tgbotapi.Update) error
}
