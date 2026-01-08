package transactionrepo

import (
	"context"
	"fmt"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/transaction"
	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

type TransactionRepository struct {
	tracker Tracker
}

func NewTransactionRepository(tracker Tracker) (ports.TransactionRepository, error) {
	if tracker == nil {
		return nil, errs.NewValueIsRequiredError("tracker")
	}

	return &TransactionRepository{tracker: tracker}, nil
}

func (t TransactionRepository) Add(ctx context.Context, tr *transaction.Transaction) error {
	stmt := `INSERT INTO transactions (id, amount, category_id, created_at, user_id)
			 VALUES ($1, $2, $3, $4, $5)`
	_, err := t.tracker.Tx().ExecContext(ctx, stmt, tr.ID(), tr.Amount(), tr.CategoryID(), tr.CreatedAt(), tr.UserID())
	if err != nil {
		return fmt.Errorf("transaction repo add: %w", err)
	}

	return nil
}

func (t TransactionRepository) Get(ctx context.Context, id shared.ID) (*transaction.Transaction, error) {
	panic("implement me")
}

func (t TransactionRepository) Update(ctx context.Context, transaction *transaction.Transaction) error {
	panic("implement me")
}

func (t TransactionRepository) Delete(ctx context.Context, id shared.ID) error {
	panic("implement me")
}
