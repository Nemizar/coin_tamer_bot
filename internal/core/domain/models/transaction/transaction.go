// Package transaction определяет доменную модель Transaction и объекты-значения.
// Обрабатывает создание финансовых транзакций, валидацию и бизнес-правила.
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

// Transaction представляет сущность финансовой транзакции с суммой, категорией и типом.
// Является основной сущностью для отслеживания доходов и расходов.
type Transaction struct {
	baseAggregate   *ddd.BaseAggregate[shared.ID]
	userID          shared.ID
	amount          Amount
	categoryID      shared.ID
	transactionType Type
	createdAt       time.Time
}

// NewTransaction создает новый экземпляр Transaction с валидированными параметрами.
// Возвращает ошибку если валидация user id или category id не проходит.
func NewTransaction(uID shared.ID, amount Amount, cID shared.ID, transactionType Type) (*Transaction, error) {
	if uID.IsZero() {
		return nil, fmt.Errorf("%w: %s", ErrInvalidUserID, uID)
	}

	if cID.IsZero() {
		return nil, fmt.Errorf("%w: %s", ErrInvalidCategoryID, cID)
	}

	return &Transaction{
		baseAggregate:   ddd.NewBaseAggregate(shared.NewID()),
		userID:          uID,
		amount:          amount,
		categoryID:      cID,
		transactionType: transactionType,
		createdAt:       time.Now(),
	}, nil
}

// CreatedAt возвращает временную метку создания транзакции.
func (t Transaction) CreatedAt() time.Time {
	return t.createdAt
}

// Type возвращает тип транзакции (доход или расход).
func (t Transaction) Type() Type {
	return t.transactionType
}

// CategoryID возвращает идентификатор категории, к которой относится транзакция.
func (t Transaction) CategoryID() shared.ID {
	return t.categoryID
}

// Amount возвращает сумму транзакции.
func (t Transaction) Amount() Amount {
	return t.amount
}

// UserID возвращает идентификатор пользователя, которому принадлежит транзакция.
func (t Transaction) UserID() shared.ID {
	return t.userID
}

// ID возвращает уникальный идентификатор транзакции.
func (t Transaction) ID() shared.ID {
	return t.baseAggregate.ID()
}

// Equals сравнивает текущую транзакцию с другой транзакцией.
// Возвращает true, если транзакции имеют одинаковый идентификатор.
func (t Transaction) Equals(other *Transaction) bool {
	if other == nil {
		return false
	}

	return t.baseAggregate.Equal(other.baseAggregate)
}
