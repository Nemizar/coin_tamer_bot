package telegram

import (
	"context"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"

	"github.com/Nemizar/coin_tamer_bot/internal/core/application/usecases/queries"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/category"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/transaction"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/user"
)

func (b *Bot) handleMsg(ctx context.Context, update tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}

	chatID := update.Message.Chat.ID
	text := strings.TrimSpace(update.Message.Text)
	if text == "" {
		return nil
	}

	u, err := b.getOrNotifyUser(ctx, chatID)
	if err != nil || u == nil {
		return err
	}

	amount, operationType, err := parseAmount(text)
	if err != nil {
		b.sendValidationError(chatID)
		return err
	}

	b.savePendingTransaction(chatID, amount)

	categories, err := b.getUserCategories(ctx, u.ID(), operationType)
	if err != nil {
		b.sendCategoriesError(chatID)
		return err
	}

	return b.sendCategoriesKeyboard(chatID, categories)
}

func (b *Bot) getOrNotifyUser(ctx context.Context, chatID int64) (*user.User, error) {
	userQuery, err := queries.NewGetUserQuery(
		strconv.FormatInt(chatID, 10),
		user.ProviderTelegram,
	)
	if err != nil {
		return nil, err
	}

	u, err := b.getUserQueryHandler.Handle(ctx, userQuery)
	if err != nil {
		return nil, err
	}

	if u == nil {
		if err = b.sendMsg(chatID, "Для начала работы необходимо зарегистрироваться /start"); err != nil {
			return nil, err
		}
	}

	return u, nil
}

func parseAmount(msg string) (transaction.Amount, category.Type, error) {
	msg = strings.TrimSpace(msg)

	opType := category.TypeExpense
	if strings.HasPrefix(msg, "+") {
		opType = category.TypeIncome
		msg = strings.TrimPrefix(msg, "+")
	}

	amount, err := transaction.NewAmountFromString(msg)
	return amount, opType, err
}

func (b *Bot) getUserCategories(
	ctx context.Context,
	userID shared.ID,
	operationType category.Type,
) ([]*category.Category, error) {
	query := queries.NewGetUserCategoriesByType(userID, operationType)

	return b.getUserCategoriesByTypeQueryHandler.Handle(ctx, query)
}

func (b *Bot) sendCategoriesKeyboard(
	chatID int64,
	categories []*category.Category,
) error {
	keyboard := newCategoriesInlineKeyboard(categories)

	if err := b.sendReplyMarkup(chatID, "Выберите категорию:", &keyboard); err != nil {
		b.logger.Error(
			"Ошибка отправки клавиатуры с категориями",
			"err", err.Error(),
		)
		return err
	}

	return nil
}
