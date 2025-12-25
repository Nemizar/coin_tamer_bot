package user

import (
	"time"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/ddd"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

type ExternalIdentity struct {
	baseAggregate *ddd.BaseEntity[shared.ID]
	userID        shared.ID
	provider      Provider
	externalID    string
	createdAt     time.Time
}

func (e ExternalIdentity) ID() shared.ID {
	return e.baseAggregate.ID()
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

func (e ExternalIdentity) GetCreatedAt() time.Time {
	return e.createdAt
}

func NewExternalIdentity(userID shared.ID, provider Provider, externalID string) (*ExternalIdentity, error) {
	if !provider.IsValid() {
		return nil, errs.NewValueIsInvalidError("provider")
	}

	if externalID == "" {
		return nil, errs.NewValueIsRequiredError("externalID")
	}

	ei := ExternalIdentity{
		baseAggregate: ddd.NewBaseEntity(shared.NewID()),
		userID:        userID,
		provider:      provider,
		externalID:    externalID,
		createdAt:     time.Now(),
	}

	return &ei, nil
}
