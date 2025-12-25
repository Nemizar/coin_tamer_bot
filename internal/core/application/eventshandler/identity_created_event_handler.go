package eventshandler

import (
	"context"
	"fmt"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/user"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/ddd"
)

type externaklIdentityAddedEventHandler struct {
}

func NewExternalIdentityAddedEventHandler() ddd.EventHandler {
	return &externaklIdentityAddedEventHandler{}
}

func (u externaklIdentityAddedEventHandler) Handle(ctx context.Context, domainEvent ddd.DomainEvent) error {
	fmt.Println("event handled")

	createEvent, ok := domainEvent.(*user.ExternalIdentityAddedEvent)
	if !ok {
		return fmt.Errorf("unexpected domain event type: %T", domainEvent)
	}

	fmt.Println(createEvent.GetID())
	fmt.Println(createEvent.GetName())
	fmt.Println(createEvent.GetPayload())

	return nil
}
