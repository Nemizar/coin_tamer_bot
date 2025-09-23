// Package user определяет доменную модель User и связанную бизнес-логику.
// Обрабатывает создание пользователей, валидацию и доменные операции.
package user

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/ddd"
)

var (
	ErrInvalidChatID = errors.New("invalid chat id")
	ErrInvalidName   = errors.New("invalid name")
)

// User представляет сущность пользователя бота с уникальной идентификацией и контекстом чата.
// Является агрегатом для пользовательских доменных операций.
type User struct {
	baseAggregate *ddd.BaseAggregate[shared.ID]
	chatID        int64
	createdAt     time.Time
	name          string
}

// NewUser создает новый экземпляр User с валидированными параметрами.
// Возвращает ошибку если валидация chat id или имени не проходит.
func NewUser(chatID int64, name string) (*User, error) {
	if chatID <= 0 {
		return nil, fmt.Errorf("chatID %d: %w", chatID, ErrInvalidChatID)
	}

	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("name %s: %w", name, ErrInvalidName)
	}

	return &User{
		baseAggregate: ddd.NewBaseAggregate(shared.NewID()),
		chatID:        chatID,
		createdAt:     time.Now(),
		name:          name,
	}, nil
}

// ID возвращает уникальный идентификатор пользователя.
func (u User) ID() shared.ID {
	return u.baseAggregate.ID()
}

// ChatID возвращает идентификатор чата пользователя в Telegram.
func (u User) ChatID() int64 {
	return u.chatID
}

// CreatedAt возвращает временную метку создания записи пользователя.
func (u User) CreatedAt() time.Time {
	return u.createdAt
}

// Name возвращает имя пользователя.
func (u User) Name() string {
	return u.name
}

// Equals сравнивает текущего пользователя с другим пользователем.
func (u User) Equals(other *User) bool {
	if other == nil {
		return false
	}

	return u.baseAggregate.Equal(other.baseAggregate)
}
