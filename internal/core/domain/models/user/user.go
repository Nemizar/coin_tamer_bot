// Package user определяет доменную модель User и связанную бизнес-логику.
// Обрабатывает создание пользователей, валидацию и доменные операции.
package user

import (
	"strings"
	"time"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/ddd"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

// User представляет сущность пользователя бота с уникальной идентификацией и контекстом чата.
// Является агрегатом для пользовательских доменных операций.
type User struct {
	baseAggregate *ddd.BaseAggregate[shared.ID]
	createdAt     time.Time
	name          string
}

// New создает новый экземпляр User с валидированными параметрами.
// Возвращает ошибку если валидация имени не проходит.
func New(name string) (*User, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errs.NewValueIsRequiredError("name")
	}

	return &User{
		baseAggregate: ddd.NewBaseAggregate(shared.NewID()),
		createdAt:     time.Now(),
		name:          name,
	}, nil
}

func Restore(id shared.ID, name string, createdAt time.Time) *User {
	return &User{
		baseAggregate: ddd.NewBaseAggregate(id),
		name:          name,
		createdAt:     createdAt,
	}
}

// ID возвращает уникальный идентификатор пользователя.
func (u User) ID() shared.ID {
	return u.baseAggregate.ID()
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
