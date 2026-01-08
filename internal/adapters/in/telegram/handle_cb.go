package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"

	"github.com/Nemizar/coin_tamer_bot/internal/core/application/usecases/commands"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
)

func (b *Bot) handleCb(ctx context.Context, cb *tgbotapi.CallbackQuery) error {
	chatID := cb.Message.Chat.ID
	prevMsgID := cb.Message.MessageID

	u, err := b.getOrNotifyUser(ctx, chatID)
	if err != nil || u == nil {
		return err
	}

	us, err := b.getUserState(chatID)
	if err != nil {
		return err
	}

	if us != UserStateWaitingForCategory {
		return err
	}

	pt, err := b.getPendingTransaction(chatID)
	if err != nil {
		return err
	}

	cID, err := uuid.Parse(cb.Data)
	if err != nil {
		return err
	}

	cmd, err := commands.NewCreateTransactionCommand(u.ID(), pt.Amount, shared.RestoreID(cID))
	if err != nil {
		return err
	}

	err = b.createTransactionCommandHandler.Handle(ctx, cmd)
	if err != nil {
		b.logger.Error(err.Error())

		err = b.sendMessageAndDeleteInlineKeyboard(chatID, prevMsgID, "Ошибка при сохранении расхода. Попробуйте еще раз")
		if err != nil {
			b.logger.Error(err.Error())
		}
	} else {
		err = b.sendMessageAndDeleteInlineKeyboard(chatID, prevMsgID, "✅ Транзакция записана!")
		if err != nil {
			b.logger.Error(err.Error())
		}
	}

	b.cache.Delete(userStateKey(chatID))
	b.cache.Delete(pendingTransactionKey(chatID))

	return nil
}
