package ports

import (
	"context"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/user"
)

type UserRepository interface {
	Create(ctx context.Context, user *user.User) error
	FindByExternalProvider(ctx context.Context, provider user.Provider, externalID string) (*user.User, error)
}
