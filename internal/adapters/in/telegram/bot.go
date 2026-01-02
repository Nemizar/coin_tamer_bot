package telegram

import (
	"context"

	"github.com/Nemizar/coin_tamer_bot/internal/core/application/usecases/commands"
	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	logger ports.Logger

	bot *tgbotapi.BotAPI

	router *router

	telegramUserRegistrationCommandHandler commands.UserRegistrationCommandHandler
}

func NewBot(
	logger ports.Logger,
	telegramBotToken string,
	telegramUserRegistrationHandler commands.UserRegistrationCommandHandler,
	createDefaultCategoriesCommandHandler commands.CreateDefaultCategoryCommandHandler,
) (*Bot, error) {
	if logger == nil {
		return nil, errs.NewValueIsRequiredError("logger")
	}

	if telegramUserRegistrationHandler == nil {
		return nil, errs.NewValueIsRequiredError("telegramUserRegistrationHandler")
	}

	if createDefaultCategoriesCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("createDefaultCategoriesCommandHandler")
	}

	if telegramBotToken == "" {
		return nil, errs.NewValueIsRequiredError("telegramBotToken")
	}

	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		return nil, err
	}

	r := newRouter()
	startCmdHandler, err := newStartCommandHandler(telegramUserRegistrationHandler)
	if err != nil {
		return nil, err
	}
	r.registerCommand(startCmdHandler)

	createDefaultCategoriesCmdHandler, err := newCreateDefaultCategoriesCommandHandler(createDefaultCategoriesCommandHandler)
	if err != nil {
		return nil, err
	}
	r.registerCommand(createDefaultCategoriesCmdHandler)

	return &Bot{
		bot:                                    bot,
		logger:                                 logger,
		router:                                 r,
		telegramUserRegistrationCommandHandler: telegramUserRegistrationHandler,
	}, nil
}

func (b *Bot) HandleUpdates(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		b.logger.Info("handle updates")

		err := b.router.handleUpdate(ctx, update)
		if err != nil {
			b.logger.Error("handle updates", "err", err)
		}
	}
}
