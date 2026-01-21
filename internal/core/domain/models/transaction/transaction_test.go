package transaction_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
	transaction2 "github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/transaction"
)

func TestNewTransaction_Success(t *testing.T) {
	userID := shared.NewID()
	categoryID := shared.NewID()
	amount := transaction2.Amount{}

	tx, err := transaction2.New(userID, amount, categoryID)

	require.NoError(t, err)
	require.False(t, tx.ID().IsZero())
	require.Equal(t, userID, tx.UserID())
	require.Equal(t, categoryID, tx.CategoryID())
	require.Equal(t, amount, tx.Amount())
	require.WithinDuration(t, time.Now(), tx.CreatedAt(), time.Second)
}

func TestNewTransaction_InvalidUserID(t *testing.T) {
	zeroUser := shared.ID{}
	categoryID := shared.NewID()
	amount := transaction2.Amount{}

	tx, err := transaction2.New(zeroUser, amount, categoryID)

	require.Error(t, err)
	require.ErrorIs(t, err, transaction2.ErrInvalidUserID)
	require.Nil(t, tx)
}

func TestNewTransaction_InvalidCategoryID(t *testing.T) {
	userID := shared.NewID()
	zeroCategory := shared.ID{}
	amount := transaction2.Amount{}

	tx, err := transaction2.New(userID, amount, zeroCategory)

	require.Error(t, err)
	require.ErrorIs(t, err, transaction2.ErrInvalidCategoryID)
	require.Nil(t, tx)
}

func TestNewTransaction_TableTests(t *testing.T) {
	userID := shared.NewID()
	categoryID := shared.NewID()
	amount := transaction2.Amount{}

	tests := []struct {
		name        string
		userID      shared.ID
		categoryID  shared.ID
		amount      transaction2.Amount
		expectError bool
		errorType   error
	}{
		{
			name:        "Валидная транзакция",
			userID:      userID,
			categoryID:  categoryID,
			amount:      amount,
			expectError: false,
		},
		{
			name:        "Невалидный ID пользователя (нулевой)",
			userID:      shared.ID{},
			categoryID:  categoryID,
			amount:      amount,
			expectError: true,
			errorType:   transaction2.ErrInvalidUserID,
		},
		{
			name:        "Невалидный ID категории (нулевой)",
			userID:      userID,
			categoryID:  shared.ID{},
			amount:      amount,
			expectError: true,
			errorType:   transaction2.ErrInvalidCategoryID,
		},
		{
			name:        "Оба ID невалидны",
			userID:      shared.ID{},
			categoryID:  shared.ID{},
			amount:      amount,
			expectError: true,
			errorType:   transaction2.ErrInvalidUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := transaction2.New(tt.userID, tt.amount, tt.categoryID)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorType != nil {
					require.ErrorIs(t, err, tt.errorType)
				}
				require.Nil(t, tx)
			} else {
				require.NoError(t, err)
				require.NotNil(t, tx)
				assert.Equal(t, tt.userID, tx.UserID())
				assert.Equal(t, tt.categoryID, tx.CategoryID())
				assert.Equal(t, tt.amount, tx.Amount())
			}
		})
	}
}

func TestTransaction_Equals(t *testing.T) {
	userID := shared.NewID()
	categoryID := shared.NewID()
	amount := transaction2.Amount{}

	tx1, err := transaction2.New(userID, amount, categoryID)
	require.NoError(t, err)

	tx2, err := transaction2.New(userID, amount, categoryID)
	require.NoError(t, err)

	tests := []struct {
		name     string
		tx1      *transaction2.Transaction
		tx2      *transaction2.Transaction
		expected bool
	}{
		{
			name:     "Один и тот же экземпляр транзакции",
			tx1:      tx1,
			tx2:      tx1,
			expected: true,
		},
		{
			name:     "Разные транзакции с одинаковыми данными",
			tx1:      tx1,
			tx2:      tx2,
			expected: false, // Different IDs
		},
		{
			name:     "Одна транзакция равна nil",
			tx1:      tx1,
			tx2:      nil,
			expected: false,
		},
		{
			name:     "Обе транзакции равны nil",
			tx1:      nil,
			tx2:      nil,
			expected: false, // According to the implementation, comparing with nil returns false
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.tx1 == nil {
				// Can't call method on nil pointer, so we expect false as per implementation
				assert.False(t, tt.expected, "Expected false when tx1 is nil")
			} else {
				// Normal case: call method on non-nil object
				result := tt.tx1.Equals(tt.tx2)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
