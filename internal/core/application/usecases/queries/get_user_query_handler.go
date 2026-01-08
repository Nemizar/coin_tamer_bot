package queries

import (
	"context"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/user"
	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

type GetUserQueryHandler interface {
	Handle(ctx context.Context, query GetUserQuery) (*user.User, error)
}

type getUserQueryHandler struct {
	uow ports.UnitOfWork
}

func NewGetUserQueryHandler(uow ports.UnitOfWork) (GetUserQueryHandler, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequiredError("uow")
	}

	return &getUserQueryHandler{uow: uow}, nil
}

func (h getUserQueryHandler) Handle(ctx context.Context, query GetUserQuery) (*user.User, error) {
	return h.uow.UserRepository().FindByExternalProvider(ctx, query.Provider(), query.ExternalID())
}
