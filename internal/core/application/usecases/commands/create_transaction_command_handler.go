package commands

import (
	"context"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/transaction"
	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

type CreateTransactionCommandHandler interface {
	Handle(ctx context.Context, command CreateTransactionCommand) error
}

type createTransactionCommandHandler struct {
	logger ports.Logger
	uow    ports.UnitOfWork
}

func NewCreateTransactionCommandHandler(logger ports.Logger, uow ports.UnitOfWork) (CreateTransactionCommandHandler, error) {
	if logger == nil {
		return nil, errs.NewValueIsRequiredError("logger")
	}

	if uow == nil {
		return nil, errs.NewValueIsRequiredError("uow")
	}

	return &createTransactionCommandHandler{
		logger: logger,
		uow:    uow,
	}, nil
}

func (t createTransactionCommandHandler) Handle(ctx context.Context, command CreateTransactionCommand) error {
	defer func(uow ports.UnitOfWork) {
		err := uow.RollbackUnlessCommitted()
		if err != nil {
			t.logger.Error("create transaction command handler: rollback failed", "err", err)
		}
	}(t.uow)

	err := t.uow.Begin(ctx)
	if err != nil {
		return err
	}

	nt, err := transaction.New(command.UserID(), command.Amount(), command.CategoryID())
	if err != nil {
		return err
	}

	err = t.uow.TransactionRepository().Add(ctx, nt)
	if err != nil {
		return err
	}

	return t.uow.Commit(ctx)
}
