package transaction_test

import (
	"testing"
	"time"

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
