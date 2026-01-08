package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) error {
	switch {
	case update.Message != nil && update.Message.IsCommand():
		cmd := update.Message.Command()
		switch cmd {
		case "start":
			return b.handleStartCommand(ctx, update)
		case "create_default_categories":
			return b.handleCreateDefaultCategoriesCommand(ctx, update)
		}

		return errs.NewValueIsInvalidErrorWithCause("command", errs.NewValueIsInvalidError("command "+cmd))
	case update.CallbackQuery != nil:
		return b.handleCb(ctx, update.CallbackQuery)
	case update.Message != nil:
		return b.handleMsg(ctx, update)
	}

	return fmt.Errorf("not expected type")
}
