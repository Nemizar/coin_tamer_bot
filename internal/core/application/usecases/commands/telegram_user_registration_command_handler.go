package commands

import (
	"context"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/identity"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/user"
	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"
)

type TelegramUserRegistrationCommandHandler interface {
	Handle(ctx context.Context, command TelegramUserRegistrationCommand) error
}

var _ TelegramUserRegistrationCommandHandler = telegramUserRegistrationCommandHandler{}

type telegramUserRegistrationCommandHandler struct {
	logger     ports.Logger
	uowFactory ports.UnitOfWorkFactory
}

func NewUserRegistrationCommandHandler(logger ports.Logger, uowFactory ports.UnitOfWorkFactory) TelegramUserRegistrationCommandHandler {
	return &telegramUserRegistrationCommandHandler{
		logger:     logger,
		uowFactory: uowFactory,
	}
}

func (u telegramUserRegistrationCommandHandler) Handle(ctx context.Context, command TelegramUserRegistrationCommand) error {
	uow, err := u.uowFactory.New(ctx)
	if err != nil {
		return err
	}

	defer func(uow ports.UnitOfWork) {
		err := uow.RollbackUnlessCommitted()
		if err != nil {
			u.logger.Error("user registration command handler: rollback failed", "err", err)
		}
	}(uow)

	nu, err := user.New(command.Username())
	if err != nil {
		return err
	}

	exIdentity, err := identity.NewExternalIdentity(nu.ID(), identity.ProviderTelegram, command.TelegramChatID())
	if err != nil {
		return err
	}

	err = uow.Begin()
	if err != nil {
		return err
	}

	err = uow.UserRepository().Create(ctx, nu)
	if err != nil {
		return err
	}

	err = uow.ExternalIdentityRepository().Add(ctx, exIdentity)
	if err != nil {
		return err
	}

	return uow.Commit()
}
