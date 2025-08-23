package transaction_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Nemizar/coin_tamer_bot/internal/domain/shared"
	"github.com/Nemizar/coin_tamer_bot/internal/domain/transaction"
)

func TestNewTransaction_Success(t *testing.T) {
	userID := shared.NewID()
	categoryID := shared.NewID()
	amount := transaction.Amount{}
	tType, err := transaction.ParseType("expense")
	require.NoError(t, err)

	tx, err := transaction.NewTransaction(userID, amount, categoryID, tType)

	require.NoError(t, err)
	require.False(t, tx.ID.IsZero())
	require.Equal(t, userID, tx.UserID)
	require.Equal(t, categoryID, tx.CategoryID)
	require.Equal(t, amount, tx.Amount)
	require.Equal(t, tType, tx.Type)
	require.WithinDuration(t, time.Now(), tx.CreatedAt, time.Second)
}

func TestNewTransaction_InvalidUserID(t *testing.T) {
	zeroUser := shared.ID{}
	categoryID := shared.NewID()
	amount := transaction.Amount{}
	tType, err := transaction.ParseType("income")
	require.NoError(t, err)

	tx, err := transaction.NewTransaction(zeroUser, amount, categoryID, tType)

	require.Error(t, err)
	require.ErrorIs(t, err, transaction.ErrInvalidUserID)
	require.True(t, tx.ID.IsZero())
}

func TestNewTransaction_InvalidCategoryID(t *testing.T) {
	userID := shared.NewID()
	zeroCategory := shared.ID{}
	amount := transaction.Amount{}
	tType, err := transaction.ParseType("income")
	require.NoError(t, err)

	tx, err := transaction.NewTransaction(userID, amount, zeroCategory, tType)

	require.Error(t, err)
	require.ErrorIs(t, err, transaction.ErrInvalidCategoryID)
	require.True(t, tx.ID.IsZero())
}
