package extidentityrepo

import (
	"context"
	"fmt"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/identity"
	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"
)

type ExternalIdentityRepository struct {
	uow Tracker
}

func NewExternalIdentityRepository(uow Tracker) ports.ExternalIdentityRepository {
	return &ExternalIdentityRepository{uow: uow}
}

func (e ExternalIdentityRepository) Add(ctx context.Context, ei *identity.ExternalIdentity) error {
	stmt := `INSERT INTO external_identities (id, user_id, provider, external_id)
				 VALUES ($1, $2, $3, $4)`
	_, err := e.uow.DB().ExecContext(ctx, stmt, ei.ID(), ei.UserID(), ei.Provider(), ei.ExternalID())
	if err != nil {
		return fmt.Errorf("external identity repo insert: %w", err)
	}

	return nil
}
