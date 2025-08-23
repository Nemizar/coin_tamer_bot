// Package user определяет доменную модель User и связанную бизнес-логику.
// Обрабатывает создание пользователей, валидацию и доменные операции.
package user

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Nemizar/coin_tamer_bot/internal/domain/shared"
)

var (
	ErrInvalidChatID = errors.New("invalid chat id")
	ErrInvalidName   = errors.New("invalid name")
)

// User представляет сущность пользователя бота с уникальной идентификацией и контекстом чата.
// Является агрегатом для пользовательских доменных операций.
type User struct {
	ID        shared.ID
	ChatID    int64
	CreatedAT time.Time
	Name      string
}

// NewUser создает новый экземпляр User с валидированными параметрами.
// Возвращает ошибку если валидация chat ID или имени не проходит.
func NewUser(chatID int64, name string) (User, error) {
	if chatID <= 0 {
		return User{}, fmt.Errorf("chatID %d: %w", chatID, ErrInvalidChatID)
	}

	if strings.TrimSpace(name) == "" {
		return User{}, fmt.Errorf("name %s: %w", name, ErrInvalidName)
	}

	return User{
		ID:        shared.NewID(),
		ChatID:    chatID,
		CreatedAT: time.Now(),
		Name:      name,
	}, nil
}
