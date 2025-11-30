package identity

import (
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/ddd"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

type ExternalIdentity struct {
	baseEntity *ddd.BaseEntity[shared.ID]
	userID     shared.ID
	provider   Provider
	externalID string
}

func (e ExternalIdentity) ID() shared.ID {
	return e.baseEntity.ID()
}

func (e ExternalIdentity) UserID() shared.ID {
	return e.userID
}

func (e ExternalIdentity) Provider() Provider {
	return e.provider
}

func (e ExternalIdentity) ExternalID() string {
	return e.externalID
}

func NewExternalIdentity(userID shared.ID, provider Provider, externalID string) (*ExternalIdentity, error) {
	if !provider.IsValid() {
		return nil, errs.NewValueIsInvalidError("provider")
	}

	if externalID == "" {
		return nil, errs.NewValueIsRequiredError("externalID")
	}

	return &ExternalIdentity{
		baseEntity: ddd.NewBaseEntity(shared.NewID()),
		userID:     userID,
		provider:   provider,
		externalID: externalID,
	}, nil
}
