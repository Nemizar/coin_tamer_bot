package commands

import (
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/user"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

type CreateDefaultCategoryCommand interface {
	ExternalID() string
	Provider() user.Provider
}

type createDefaultCategoryCommand struct {
	externalID string
	provider   user.Provider
}

func (c createDefaultCategoryCommand) ExternalID() string {
	return c.externalID
}

func (c createDefaultCategoryCommand) Provider() user.Provider {
	return c.provider
}

func NewCreateDefaultCategoryCommand(externalID string, provider user.Provider) (CreateDefaultCategoryCommand, error) {
	if externalID == "" || externalID == "0" {
		return nil, errs.NewValueIsRequiredError("externalID")
	}

	if !provider.IsValid() {
		return nil, errs.NewValueIsRequiredError("provider")
	}

	return createDefaultCategoryCommand{externalID: externalID, provider: provider}, nil
}
