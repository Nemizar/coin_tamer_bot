package telegram

import (
	"context"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/user"

	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"

	"github.com/Nemizar/coin_tamer_bot/internal/core/application/usecases/commands"
)

type createDefaultCategoriesCommandHandler struct {
	usecase commands.CreateDefaultCategoryCommandHandler
}

func newCreateDefaultCategoriesCommandHandler(usecase commands.CreateDefaultCategoryCommandHandler) (*createDefaultCategoriesCommandHandler, error) {
	if usecase == nil {
		return nil, errs.NewValueIsRequiredError("usecase")
	}

	return &createDefaultCategoriesCommandHandler{usecase: usecase}, nil
}

func (h *createDefaultCategoriesCommandHandler) command() string {
	return "create_default_categories"
}

func (h *createDefaultCategoriesCommandHandler) handle(ctx context.Context, update tgbotapi.Update) error {
	cmd, err := commands.NewCreateDefaultCategoryCommand(
		strconv.FormatInt(update.Message.Chat.ID, 10),
		user.ProviderTelegram,
	)

	if err != nil {
		return err
	}

	return h.usecase.Handle(ctx, cmd)
}
