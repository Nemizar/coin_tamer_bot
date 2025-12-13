package commands

import (
	"strings"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/identity"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

type UserRegistrationCommand interface {
	Username() string
	ChatID() string
	Provider() identity.Provider
}

type userRegistrationCommand struct {
	username string
	chatID   string
	provider identity.Provider
}

func (u userRegistrationCommand) Username() string {
	return u.username
}

func (u userRegistrationCommand) ChatID() string {
	return u.chatID
}

func (u userRegistrationCommand) Provider() identity.Provider {
	return u.provider
}

func NewUserRegistrationCommand(username string, chatID string, provider identity.Provider) (UserRegistrationCommand, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, errs.NewValueIsRequiredError("username")
	}

	if chatID == "0" || chatID == "" {
		return nil, errs.NewValueIsRequiredError("chatID")
	}

	tID := chatID

	return &userRegistrationCommand{username: username, chatID: tID, provider: provider}, nil
}
