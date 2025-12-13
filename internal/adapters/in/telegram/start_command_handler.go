package telegram

import (
	"context"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"

	"github.com/Nemizar/coin_tamer_bot/internal/core/application/usecases/commands"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/identity"
)

type startCommandHandler struct {
	usecase commands.UserRegistrationCommandHandler
}

func newStartCommandHandler(usecase commands.UserRegistrationCommandHandler) (*startCommandHandler, error) {
	if usecase == nil {
		return nil, errs.NewValueIsRequiredError("usecase")
	}

	return &startCommandHandler{usecase: usecase}, nil
}

func (h *startCommandHandler) command() string {
	return "start"
}

func (h *startCommandHandler) handle(ctx context.Context, update tgbotapi.Update) error {
	cmd, err := commands.NewUserRegistrationCommand(
		update.Message.From.UserName,
		strconv.FormatInt(update.Message.Chat.ID, 10),
		identity.ProviderTelegram,
	)

	if err != nil {
		return err
	}

	return h.usecase.Handle(ctx, cmd)
}
