package commands

import (
	"fmt"
	"strings"

	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

type TelegramUserRegistrationCommand interface {
	Username() string
	TelegramChatID() string
}

type telegramUserRegistrationCommand struct {
	username       string
	telegramChatID string
}

func (u telegramUserRegistrationCommand) Username() string {
	return u.username
}

func (u telegramUserRegistrationCommand) TelegramChatID() string {
	return u.telegramChatID
}

func NewUserRegistrationCommand(username string, telegramChatID int) (TelegramUserRegistrationCommand, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, errs.NewValueIsRequiredError("username")
	}

	if telegramChatID == 0 {
		return nil, errs.NewValueIsRequiredError("telegramChatID")
	}

	tID := fmt.Sprintf("%d", telegramChatID)

	return &telegramUserRegistrationCommand{username: username, telegramChatID: tID}, nil
}
