package transaction_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/transaction"
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
		{"ноль", decimal.NewFromInt(0)},
		{"отрицательное", decimal.NewFromInt(-5)},
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
		{"3,14", "3.14"},                 // поддержка запятой
		{"0.01", "0.01"},                 // минимальная положительная сумма
		{"999999999.99", "999999999.99"}, // большая сумма
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

func TestNewAmountFromString_WithLeadingZeros(t *testing.T) {
	amount, err := transaction.NewAmountFromString("000123.45")
	require.NoError(t, err)
	require.Equal(t, "123.45", amount.String())
}

func TestNewAmountFromString_WithTrailingSpaces(t *testing.T) {
	amount, err := transaction.NewAmountFromString("  123.45  ")
	require.NoError(t, err)
	require.Equal(t, "123.45", amount.String())
}

func TestNewAmountFromString_WithVerySmallPositiveValue(t *testing.T) {
	amount, err := transaction.NewAmountFromString("0.001")
	require.NoError(t, err)
	require.Equal(t, "0.00", amount.String()) // Should round to 2 decimal places
}

func TestAmount_ValueAndString(t *testing.T) {
	dec := decimal.NewFromFloat(123.45)
	amount, err := transaction.NewAmount(dec)
	require.NoError(t, err)

	assert.Equal(t, dec, amount.Value())
	assert.Equal(t, "123.45", amount.String())
}

func TestAmount_Equals(t *testing.T) {
	dec1 := decimal.NewFromFloat(123.45)
	dec2 := decimal.NewFromFloat(123.45)
	dec3 := decimal.NewFromFloat(99.99)

	amount1, err := transaction.NewAmount(dec1)
	require.NoError(t, err)

	amount2, err := transaction.NewAmount(dec2)
	require.NoError(t, err)

	amount3, err := transaction.NewAmount(dec3)
	require.NoError(t, err)

	assert.Equal(t, amount1, amount2)    // Same value
	assert.NotEqual(t, amount1, amount3) // Different values
}
