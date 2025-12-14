package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Nemizar/coin_tamer_bot/internal/adapters/out/postgres/extidentityrepo"

	"github.com/Nemizar/coin_tamer_bot/internal/adapters/out/postgres/categoryrepo"
	"github.com/Nemizar/coin_tamer_bot/internal/adapters/out/postgres/transactionrepo"
	"github.com/Nemizar/coin_tamer_bot/internal/adapters/out/postgres/userrepo"
	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/ddd"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

var _ ports.UnitOfWork = &UnitOfWork{}

type UnitOfWork struct {
	pool              *sqlx.DB
	tx                *sqlx.Tx
	committed         bool
	trackedAggregates []ddd.AggregateRoot
	mediatr           ddd.Mediatr
	logger            ports.Logger

	// Ленивая инициализация репозиториев
	categoryRepo         ports.CategoryRepository
	transactionRepo      ports.TransactionRepository
	userRepo             ports.UserRepository
	externalIdentityRepo ports.ExternalIdentityRepository
}

func NewUnitOfWork(pool *sqlx.DB, mediatr ddd.Mediatr, logger ports.Logger) (ports.UnitOfWork, error) {
	if pool == nil {
		return nil, errs.NewValueIsRequiredError("pool")
	}

	if mediatr == nil {
		return nil, errs.NewValueIsRequiredError("mediatr")
	}

	if logger == nil {
		return nil, errs.NewValueIsRequiredError("logger")
	}

	return &UnitOfWork{pool: pool, mediatr: mediatr, logger: logger}, nil
}

func (u *UnitOfWork) Begin(ctx context.Context) error {
	if u.tx != nil {
		return fmt.Errorf("transaction already started")
	}

	tx, err := u.pool.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	u.tx = tx
	u.committed = false
	return nil
}

func (u *UnitOfWork) Commit(ctx context.Context) error {
	if u.tx == nil {
		return errs.NewValueIsRequiredError("cannot commit without transaction")
	}
	if u.committed {
		return fmt.Errorf("transaction already committed")
	}

	if err := u.publishDomainEvents(ctx); err != nil {
		return err
	}

	err := u.tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	u.committed = true
	u.clearTx()
	return nil
}

func (u *UnitOfWork) RollbackUnlessCommitted() error {
	if u.tx == nil || u.committed {
		return nil // nothing to rollback
	}

	err := u.tx.Rollback()
	u.clearTx()

	if err != nil && !errors.Is(err, sql.ErrTxDone) {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}

	return nil
}

func (u *UnitOfWork) Tx() *sqlx.Tx {
	return u.tx
}

func (u *UnitOfWork) InTx() bool {
	return u.tx != nil
}

func (u *UnitOfWork) DB() *sqlx.DB {
	return u.pool
}

func (u *UnitOfWork) Track(agg ddd.AggregateRoot) {
	u.trackedAggregates = append(u.trackedAggregates, agg)
}

func (u *UnitOfWork) TrackedAggregates() []ddd.AggregateRoot {
	return u.trackedAggregates
}

func (u *UnitOfWork) clearTx() {
	u.tx = nil
	u.committed = false
	// НЕ трогаем trackedAggregates — их очистка должна быть после обработки событий!
}

// --- Репозитории (ленивая инициализация) ---

func (u *UnitOfWork) CategoryRepository() ports.CategoryRepository {
	if u.categoryRepo == nil {
		u.categoryRepo = categoryrepo.NewCategoryRepository(u)
	}

	return u.categoryRepo
}

func (u *UnitOfWork) UserRepository() ports.UserRepository {
	if u.userRepo == nil {
		u.userRepo = userrepo.NewUserRepository(u)
	}

	return u.userRepo
}

func (u *UnitOfWork) TransactionRepository() ports.TransactionRepository {
	if u.transactionRepo == nil {
		u.transactionRepo = transactionrepo.NewTransactionRepository(u)
	}

	return u.transactionRepo
}

func (u *UnitOfWork) ExternalIdentityRepository() ports.ExternalIdentityRepository {
	if u.externalIdentityRepo == nil {
		u.externalIdentityRepo = extidentityrepo.NewExternalIdentityRepository(u)
	}

	return u.externalIdentityRepo
}

func (u *UnitOfWork) publishDomainEvents(ctx context.Context) error {
	for _, aggregate := range u.trackedAggregates {
		for _, event := range aggregate.GetDomainEvents() {
			err := u.mediatr.Publish(ctx, event)
			if err != nil {
				u.logger.Error("publish event", "err", err)

				continue
			}
		}

		aggregate.ClearDomainEvents()
	}

	return nil
}
