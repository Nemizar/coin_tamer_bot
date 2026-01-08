package transaction

import (
	"errors"
	"fmt"
	"time"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/ddd"
)

var (
	ErrInvalidUserID     = errors.New("invalid user id")
	ErrInvalidCategoryID = errors.New("invalid category id")
)

type Transaction struct {
	baseAggregate *ddd.BaseAggregate[shared.ID]
	userID        shared.ID
	amount        Amount
	categoryID    shared.ID
	createdAt     time.Time
}

func New(uID shared.ID, amount Amount, cID shared.ID) (*Transaction, error) {
	if uID.IsZero() {
		return nil, fmt.Errorf("%w: %s", ErrInvalidUserID, uID)
	}

	if cID.IsZero() {
		return nil, fmt.Errorf("%w: %s", ErrInvalidCategoryID, cID)
	}

	return &Transaction{
		baseAggregate: ddd.NewBaseAggregate(shared.NewID()),
		userID:        uID,
		amount:        amount,
		categoryID:    cID,
		createdAt:     time.Now(),
	}, nil
}

func (t Transaction) CreatedAt() time.Time {
	return t.createdAt
}

func (t Transaction) CategoryID() shared.ID {
	return t.categoryID
}

func (t Transaction) Amount() Amount {
	return t.amount
}

func (t Transaction) UserID() shared.ID {
	return t.userID
}

func (t Transaction) ID() shared.ID {
	return t.baseAggregate.ID()
}

func (t Transaction) Equals(other *Transaction) bool {
	if other == nil {
		return false
	}

	return t.baseAggregate.Equal(other.baseAggregate)
}
