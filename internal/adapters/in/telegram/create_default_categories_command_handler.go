package telegram

import (
	"context"
	"errors"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/user"

	"github.com/Nemizar/coin_tamer_bot/internal/core/application/usecases/commands"
)

func (b *Bot) handleCreateDefaultCategoriesCommand(ctx context.Context, update tgbotapi.Update) error {
	cmd, err := commands.NewCreateDefaultCategoryCommand(
		strconv.FormatInt(update.Message.Chat.ID, 10),
		user.ProviderTelegram,
	)
	if err != nil {
		return err
	}

	err = b.createDefaultCategoriesCommandHandler.Handle(ctx, cmd)
	if err != nil {
		var entityAlreadyExistsError *errs.EntityAlreadyExistsError

		if errors.As(err, &entityAlreadyExistsError) {
			err2 := b.sendMsg(update.Message.Chat.ID, "Категории уже созданы. Можно вести учет доходов и расходов.")
			if err2 != nil {
				b.logger.Error("Ошибка отправки сообщения о существующих категориях", err2, err2.Error())
			}

			return nil
		}

		err2 := b.sendMsg(update.Message.Chat.ID, "Ошибка при создании категорий. Попробуйте снова")
		if err2 != nil {
			b.logger.Error("Ошибка создания категорий", err2, err2.Error())
		}

		return err
	}

	err = b.sendMsg(update.Message.Chat.ID, "Категории успешно созданы. Для начала ведения расходов отправьте сумму расхода в чат")
	if err != nil {
		b.logger.Error("Ошибка отправки сообщения о создании категорий", err, err.Error())
	}

	return nil
}
