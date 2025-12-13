package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

type router struct {
	commandsHandler map[string]commandHandler
	messageHandler  []messageHandler
	callbackHandler []callbackHandler
}

func newRouter() *router {
	return &router{
		commandsHandler: make(map[string]commandHandler),
		messageHandler:  make([]messageHandler, 0),
		callbackHandler: make([]callbackHandler, 0),
	}
}

func (r *router) registerCommand(h commandHandler) {
	r.commandsHandler[h.command()] = h
}

func (r *router) handleUpdate(ctx context.Context, update tgbotapi.Update) error {
	switch {
	case update.Message != nil && update.Message.IsCommand():
		cmd := update.Message.Command()
		if handler, ok := r.commandsHandler[cmd]; ok {
			return handler.handle(ctx, update)
		}

		return errs.NewValueIsInvalidErrorWithCause("command", errs.NewValueIsInvalidError("command "+cmd))
	case update.CallbackQuery != nil:
		// TODO: implement callback query handling
	case update.Message != nil:
		// TODO: implement message handling
	}

	return fmt.Errorf("not expected type")
}
