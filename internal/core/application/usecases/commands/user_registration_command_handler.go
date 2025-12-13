package commands

import (
	"context"
	"errors"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/identity"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/user"
	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

type UserRegistrationCommandHandler interface {
	Handle(ctx context.Context, command UserRegistrationCommand) error
}

var _ UserRegistrationCommandHandler = userRegistrationCommandHandler{}

type userRegistrationCommandHandler struct {
	logger     ports.Logger
	uowFactory ports.UnitOfWorkFactory
}

func NewUserRegistrationCommandHandler(logger ports.Logger, uowFactory ports.UnitOfWorkFactory) UserRegistrationCommandHandler {
	return &userRegistrationCommandHandler{
		logger:     logger,
		uowFactory: uowFactory,
	}
}

func (u userRegistrationCommandHandler) Handle(ctx context.Context, command UserRegistrationCommand) error {
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

	existsUser, err := uow.UserRepository().FindByExternalProvider(command.Provider(), command.ChatID())
	if err != nil {
		if !errors.Is(err, errs.ErrObjectNotFound) {
			return err
		}
	}

	if existsUser != nil {
		return nil
	}

	nu, err := user.New(command.Username())
	if err != nil {
		return err
	}

	exIdentity, err := identity.NewExternalIdentity(nu.ID(), command.Provider(), command.ChatID())
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
