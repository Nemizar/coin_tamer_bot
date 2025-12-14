package transactionrepo

import (
	"context"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/transaction"
	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"
)

type TransactionRepository struct {
	tracker Tracker
}

func NewTransactionRepository(tracker Tracker) ports.TransactionRepository {
	return &TransactionRepository{tracker: tracker}
}

func (t TransactionRepository) Add(ctx context.Context, transaction *transaction.Transaction) error {
	panic("implement me")
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
