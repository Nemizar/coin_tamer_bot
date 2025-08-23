// Package transaction определяет доменную модель Transaction и объекты-значения.
// Обрабатывает создание финансовых транзакций, валидацию и бизнес-правила.
package transaction

import (
	"errors"
	"fmt"
	"time"

	"github.com/Nemizar/coin_tamer_bot/internal/domain/shared"
)

var (
	ErrInvalidUserID     = errors.New("invalid user id")
	ErrInvalidCategoryID = errors.New("invalid category id")
)

// Transaction представляет сущность финансовой транзакции с суммой, категорией и типом.
// Является основной сущностью для отслеживания доходов и расходов.
type Transaction struct {
	ID         shared.ID
	UserID     shared.ID
	Amount     Amount
	CategoryID shared.ID
	Type       Type
	CreatedAt  time.Time
}

// NewTransaction создает новый экземпляр Transaction с валидированными параметрами.
// Возвращает ошибку если валидация user ID или category ID не проходит.
func NewTransaction(uID shared.ID, amount Amount, cID shared.ID, transactionType Type) (Transaction, error) {
	if uID.IsZero() {
		return Transaction{}, fmt.Errorf("%w: %s", ErrInvalidUserID, uID)
	}

	if cID.IsZero() {
		return Transaction{}, fmt.Errorf("%w: %s", ErrInvalidCategoryID, cID)
	}

	return Transaction{
		ID:         shared.NewID(),
		UserID:     uID,
		Amount:     amount,
		CategoryID: cID,
		Type:       transactionType,
		CreatedAt:  time.Now(),
	}, nil
}
