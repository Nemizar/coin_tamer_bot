package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

type unitOfWorkFactory struct {
	db *sqlx.DB
}

func NewUnitOfWorkFactory(db *sqlx.DB) (ports.UnitOfWorkFactory, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}
	return &unitOfWorkFactory{db: db}, nil
}

func (f *unitOfWorkFactory) New(ctx context.Context) (ports.UnitOfWork, error) {
	return NewUnitOfWork(f.db)
}
