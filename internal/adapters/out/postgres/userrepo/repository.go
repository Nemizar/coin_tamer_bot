package userrepo

import (
	"context"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/user"
	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"
)

type UserRepository struct {
	uow Tracker
}

func NewUserRepository(uow Tracker) ports.UserRepository {
	return &UserRepository{uow: uow}
}

func (u UserRepository) Create(ctx context.Context, user *user.User) error {
	panic("implement me")
}
