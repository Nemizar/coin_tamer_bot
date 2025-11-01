package postgres

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"

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

	// Ленивая инициализация репозиториев
	categoryRepo    ports.CategoryRepository
	transactionRepo ports.TransactionRepository
	userRepo        ports.UserRepository
}

func NewUnitOfWork(pool *sqlx.DB) (ports.UnitOfWork, error) {
	if pool == nil {
		return nil, errs.NewValueIsRequiredError("pool")
	}

	return &UnitOfWork{pool: pool}, nil
}

// Begin начинает транзакцию. Возвращает ошибку, если не удалось.
func (u *UnitOfWork) Begin() error {
	if u.tx != nil {
		return fmt.Errorf("transaction already started")
	}

	tx, err := u.pool.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	u.tx = tx
	u.committed = false
	return nil
}

// Commit фиксирует транзакцию.
func (u *UnitOfWork) Commit() error {
	if u.tx == nil {
		return errs.NewValueIsRequiredError("cannot commit without transaction")
	}
	if u.committed {
		return fmt.Errorf("transaction already committed")
	}

	err := u.tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	u.committed = true
	u.clearTx()
	return nil
}

// RollbackUnlessCommitted откатывает, если не был сделан коммит.
func (u *UnitOfWork) RollbackUnlessCommitted() error {
	if u.tx == nil || u.committed {
		return nil // nothing to rollback
	}

	err := u.tx.Rollback()
	u.clearTx()

	if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}

	return nil
}

// Tx возвращает текущую транзакцию для репозиториев.
func (u *UnitOfWork) Tx() *sqlx.Tx {
	return u.tx
}

// InTx возвращает true, если транзакция активна.
func (u *UnitOfWork) InTx() bool {
	return u.tx != nil
}

// DB возвращает соединение (для чтения вне транзакции).
func (u *UnitOfWork) DB() *sqlx.DB {
	return u.pool
}

// Track добавляет агрегат для отслеживания (например, для событий).
func (u *UnitOfWork) Track(agg ddd.AggregateRoot) {
	u.trackedAggregates = append(u.trackedAggregates, agg)
}

// TrackedAggregates возвращает отслеживаемые агрегаты (например, для публикации событий ПОСЛЕ коммита).
func (u *UnitOfWork) TrackedAggregates() []ddd.AggregateRoot {
	return u.trackedAggregates
}

// clearTx сбрасывает транзакцию, но НЕ очищает trackedAggregates — это должно быть отдельно!
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
