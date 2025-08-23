// Package transaction определяет объект-значение Amount для денежных значений.
package transaction

import (
	"errors"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

var (
	ErrInvalidAmount       = errors.New("amount must be positive or more than zero")
	ErrInvalidAmountFormat = errors.New("invalid amount format")
)

// Amount представляет денежное значение с точной десятичной арифметикой.
// Гарантирует валидные положительные суммы и предоставляет финансовые операции.
type Amount struct {
	value decimal.Decimal
}

// NewAmountFromString создает Amount из строкового представления.
// Поддерживает как запятую, так и точку в качестве десятичного разделителя.
func NewAmountFromString(amountStr string) (Amount, error) {
	cleaned := strings.ReplaceAll(strings.TrimSpace(amountStr), ",", ".")

	value, err := decimal.NewFromString(cleaned)
	if err != nil {
		return Amount{}, fmt.Errorf("amountStr %s: %w", amountStr, ErrInvalidAmountFormat)
	}

	return NewAmount(value)
}

// NewAmountFromFloat создает Amount из значения float64.
// Примечание: применяются ограничения точности чисел с плавающей точкой.
func NewAmountFromFloat(amount float64) (Amount, error) {
	value := decimal.NewFromFloat(amount)

	return NewAmount(value)
}

// NewAmount создает Amount из значения decimal.Decimal.
// Валидирует, что сумма положительная и не нулевая.
func NewAmount(value decimal.Decimal) (Amount, error) {
	if value.IsZero() || value.IsNegative() {
		return Amount{}, ErrInvalidAmount
	}

	return Amount{value: value}, nil
}

// Value возвращает underlying decimal значение.
func (a Amount) Value() decimal.Decimal {
	return a.value
}

// String возвращает форматированное строковое представление с 2 десятичными знаками.
func (a Amount) String() string {
	return a.value.StringFixed(2)
}
