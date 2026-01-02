package eventshandler

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/user"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/ddd"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

type externalIdentityAddedEventHandler struct {
	bot *tgbotapi.BotAPI
}

func NewExternalIdentityAddedEventHandler(bot *tgbotapi.BotAPI) (ddd.EventHandler, error) {
	if bot == nil {
		return nil, errs.NewValueIsRequiredError("bot")
	}
	return &externalIdentityAddedEventHandler{
		bot: bot,
	}, nil
}

func (u externalIdentityAddedEventHandler) Handle(ctx context.Context, domainEvent ddd.DomainEvent) error {
	createEvent, ok := domainEvent.(*user.RegisterEvent)
	if !ok {
		return fmt.Errorf("unexpected domain event type: %T", domainEvent)
	}

	externalIdentity, err := strconv.Atoi(createEvent.GetPayload().ExternalID)
	if err != nil {
		return err
	}

	_, err = u.bot.Send(tgbotapi.NewMessage(int64(externalIdentity), "Hello, "+createEvent.GetPayload().UserName+"!"))
	if err != nil {
		return err
	}

	return nil
}
