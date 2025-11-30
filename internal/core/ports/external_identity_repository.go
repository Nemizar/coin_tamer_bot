package ports

import (
	"context"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/identity"
)

type ExternalIdentityRepository interface {
	Add(ctx context.Context, ei *identity.ExternalIdentity) error
}
