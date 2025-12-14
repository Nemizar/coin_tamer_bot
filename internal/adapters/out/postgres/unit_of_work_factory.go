package postgres

import (
	"github.com/jmoiron/sqlx"

	"github.com/Nemizar/coin_tamer_bot/internal/pkg/ddd"

	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

type unitOfWorkFactory struct {
	db      *sqlx.DB
	mediatr ddd.Mediatr
	logger  ports.Logger
}

func NewUnitOfWorkFactory(db *sqlx.DB, mediatr ddd.Mediatr, logger ports.Logger) (ports.UnitOfWorkFactory, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}

	if mediatr == nil {
		return nil, errs.NewValueIsRequiredError("mediatr")
	}

	if logger == nil {
		return nil, errs.NewValueIsRequiredError("logger")
	}

	return &unitOfWorkFactory{db: db, mediatr: mediatr, logger: logger}, nil
}

func (f *unitOfWorkFactory) New() (ports.UnitOfWork, error) {
	return NewUnitOfWork(f.db, f.mediatr, f.logger)
}
