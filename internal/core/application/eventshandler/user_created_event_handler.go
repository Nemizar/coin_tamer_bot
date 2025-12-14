package eventshandler

import (
	"context"
	"fmt"

	"github.com/Nemizar/coin_tamer_bot/internal/pkg/ddd"
)

type userCreatedEventHandler struct {
}

func NewUserCreatedEventHandler() ddd.EventHandler {
	return &userCreatedEventHandler{}
}

func (u userCreatedEventHandler) Handle(ctx context.Context, event ddd.DomainEvent) error {
	fmt.Println("event handled")

	return nil
}
