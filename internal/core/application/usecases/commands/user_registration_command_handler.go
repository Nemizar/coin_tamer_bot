package commands

import (
	"context"
	"errors"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/user"
	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

type UserRegistrationCommandHandler interface {
	Handle(ctx context.Context, command UserRegistrationCommand) error
}

var _ UserRegistrationCommandHandler = userRegistrationCommandHandler{}

type userRegistrationCommandHandler struct {
	logger ports.Logger
	uow    ports.UnitOfWork
}

func NewUserRegistrationCommandHandler(logger ports.Logger, uow ports.UnitOfWork) (UserRegistrationCommandHandler, error) {
	if logger == nil {
		return nil, errs.NewValueIsRequiredError("logger")
	}

	if uow == nil {
		return nil, errs.NewValueIsRequiredError("uow")
	}

	return &userRegistrationCommandHandler{
		logger: logger,
		uow:    uow,
	}, nil
}

func (u userRegistrationCommandHandler) Handle(ctx context.Context, command UserRegistrationCommand) error {
	defer func(uow ports.UnitOfWork) {
		err := uow.RollbackUnlessCommitted()
		if err != nil {
			u.logger.Error("user registration command handler: rollback failed", "err", err)
		}
	}(u.uow)

	err := u.uow.Begin(ctx)
	if err != nil {
		return err
	}

	existsUser, err := u.uow.UserRepository().FindByExternalProvider(ctx, command.Provider(), command.ChatID())
	if err != nil {
		if !errors.Is(err, errs.ErrObjectNotFound) {
			return err
		}
	}

	if existsUser != nil {
		return nil
	}

	nu, err := user.New(command.Username(), command.ChatID(), command.Provider())
	if err != nil {
		return err
	}

	err = u.uow.UserRepository().Create(ctx, nu)
	if err != nil {
		return err
	}

	return u.uow.Commit(ctx)
}
