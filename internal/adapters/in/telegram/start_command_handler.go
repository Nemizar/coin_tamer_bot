package telegram

import (
	"context"
	"errors"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Nemizar/coin_tamer_bot/internal/core/application/usecases/commands"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/user"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

func (b *Bot) handleStartCommand(ctx context.Context, update tgbotapi.Update) error {
	cmd, err := commands.NewUserRegistrationCommand(
		update.Message.From.UserName,
		strconv.FormatInt(update.Message.Chat.ID, 10),
		user.ProviderTelegram,
	)

	if err != nil {
		return err
	}

	err = b.userRegistrationCommandHandler.Handle(ctx, cmd)
	if err != nil {
		var entityAlreadyExistsError *errs.EntityAlreadyExistsError

		if errors.As(err, &entityAlreadyExistsError) {
			err2 := b.sendMsg(update.Message.Chat.ID, "Вы уже зарегистрированы. Команда /start предназначена для новых пользователей.")
			if err2 != nil {
				b.logger.Error("Ошибка отправки сообщения о повторной регистрации", err2, err2.Error())
			}
			return nil
		}

		err2 := b.sendMsg(update.Message.Chat.ID, "Ошибка регистрации. Попробуйте снова /start")
		if err2 != nil {
			b.logger.Error("Ошибка отправки сообщения об ошибке регистрации", err2, err2.Error())
		}

		return err
	}

	err = b.sendMsg(update.Message.Chat.ID, "Успешная регистрация. Выполните команду /create_default_categories для создания категорий")
	if err != nil {
		b.logger.Error("Ошибка отправки сообщения об успешной регистрации", err, err.Error())
	}

	return nil
}
