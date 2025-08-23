package transaction_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/Nemizar/coin_tamer_bot/internal/domain/transaction"
)

func TestNewAmount_Success(t *testing.T) {
	val := decimal.NewFromInt(100)
	amount, err := transaction.NewAmount(val)

	require.NoError(t, err)
	require.Equal(t, val, amount.Value())
	require.Equal(t, "100.00", amount.String())
}

func TestNewAmount_ZeroOrNegative(t *testing.T) {
	tests := []struct {
		name  string
		input decimal.Decimal
	}{
		{"zero", decimal.NewFromInt(0)},
		{"negative", decimal.NewFromInt(-5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := transaction.NewAmount(tt.input)
			require.ErrorIs(t, err, transaction.ErrInvalidAmount)
		})
	}
}

func TestNewAmountFromFloat_Success(t *testing.T) {
	amount, err := transaction.NewAmountFromFloat(12.34)
	require.NoError(t, err)

	require.Equal(t, "12.34", amount.Value().String())
	require.Equal(t, "12.34", amount.String())
}

func TestNewAmountFromFloat_Invalid(t *testing.T) {
	_, err := transaction.NewAmountFromFloat(0)
	require.ErrorIs(t, err, transaction.ErrInvalidAmount)

	_, err = transaction.NewAmountFromFloat(-3.5)
	require.ErrorIs(t, err, transaction.ErrInvalidAmount)
}

func TestNewAmountFromString_Success(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"10", "10.00"},
		{"  25.5 ", "25.50"},
		{"3,14", "3.14"}, // поддержка запятой
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			amount, err := transaction.NewAmountFromString(tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.expected, amount.String())
		})
	}
}

func TestNewAmountFromString_InvalidFormat(t *testing.T) {
	_, err := transaction.NewAmountFromString("abc")
	require.ErrorIs(t, err, transaction.ErrInvalidAmountFormat)
}

func TestNewAmountFromString_ZeroOrNegative(t *testing.T) {
	_, err := transaction.NewAmountFromString("0")
	require.ErrorIs(t, err, transaction.ErrInvalidAmount)

	_, err = transaction.NewAmountFromString("-42")
	require.ErrorIs(t, err, transaction.ErrInvalidAmount)
}
