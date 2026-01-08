package queries

import (
	"strings"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/user"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

type GetUserQuery interface {
	ExternalID() string
	Provider() user.Provider
}

type getUserQuery struct {
	externalID string
	provider   user.Provider
}

func (q getUserQuery) ExternalID() string {
	return q.externalID
}

func (q getUserQuery) Provider() user.Provider {
	return q.provider
}

func NewGetUserQuery(externalID string, provider user.Provider) (GetUserQuery, error) {
	externalID = strings.TrimSpace(externalID)

	if externalID == "" {
		return nil, errs.NewValueIsRequiredError("externalID")
	}

	return &getUserQuery{
		externalID: externalID,
		provider:   provider,
	}, nil
}
