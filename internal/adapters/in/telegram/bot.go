package telegram

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/Nemizar/coin_tamer_bot/internal/core/application/usecases/commands"
	"github.com/Nemizar/coin_tamer_bot/internal/core/application/usecases/queries"
	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	logger ports.Logger

	bot *tgbotapi.BotAPI

	cache *cache.Cache

	userRegistrationCommandHandler        commands.UserRegistrationCommandHandler
	createDefaultCategoriesCommandHandler commands.CreateDefaultCategoryCommandHandler
	createTransactionCommandHandler       commands.CreateTransactionCommandHandler

	getUserQueryHandler                 queries.GetUserQueryHandler
	getUserCategoriesByTypeQueryHandler queries.GetUserCategoriesByTypeQueryHandler
}

func NewBot(
	logger ports.Logger,
	telegramBotToken string,
	userRegistrationHandler commands.UserRegistrationCommandHandler,
	createDefaultCategoriesCommandHandler commands.CreateDefaultCategoryCommandHandler,
	createTransactionCommandHandler commands.CreateTransactionCommandHandler,
	getUserCategoriesByTypeQueryHandler queries.GetUserCategoriesByTypeQueryHandler,
	getUserQueryHandler queries.GetUserQueryHandler,
) (*Bot, error) {
	if logger == nil {
		return nil, errs.NewValueIsRequiredError("logger")
	}

	if userRegistrationHandler == nil {
		return nil, errs.NewValueIsRequiredError("userRegistrationHandler")
	}

	if createDefaultCategoriesCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("createDefaultCategoriesCommandHandler")
	}

	if createTransactionCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("createTransactionCommandHandler")
	}

	if getUserCategoriesByTypeQueryHandler == nil {
		return nil, errs.NewValueIsRequiredError("getUserCategoriesByTypeQueryHandler")
	}

	if getUserQueryHandler == nil {
		return nil, errs.NewValueIsRequiredError("getUserQueryHandler")
	}

	if telegramBotToken == "" {
		return nil, errs.NewValueIsRequiredError("telegramBotToken")
	}

	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		return nil, err
	}

	tgBot := &Bot{
		bot:                                   bot,
		logger:                                logger,
		userRegistrationCommandHandler:        userRegistrationHandler,
		createDefaultCategoriesCommandHandler: createDefaultCategoriesCommandHandler,
		createTransactionCommandHandler:       createTransactionCommandHandler,
		getUserCategoriesByTypeQueryHandler:   getUserCategoriesByTypeQueryHandler,
		getUserQueryHandler:                   getUserQueryHandler,
		cache:                                 cache.New(5*time.Minute, 10*time.Minute),
	}

	return tgBot, nil
}

func (b *Bot) HandleUpdates(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		b.logger.Info("handle update", "update_id", update.UpdateID)

		if err := b.safeHandleUpdate(ctx, update); err != nil {
			b.logger.Error(
				"failed to handle update",
				"update_id", update.UpdateID,
				"err", err.Error(),
			)
		}
	}
}

func (b *Bot) safeHandleUpdate(ctx context.Context, update tgbotapi.Update) (err error) {
	defer func() {
		if r := recover(); r != nil {
			b.logger.Error(
				"panic while handling telegram update",
				"panic", r,
				"update_id", update.UpdateID,
			)
		}
	}()

	return b.handleUpdate(ctx, update)
}
