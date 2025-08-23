// Package category определяет доменную модель Category для организации транзакций.
// Обрабатывает создание категорий, валидацию и иерархические отношения.
package category

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Nemizar/coin_tamer_bot/internal/domain/shared"
)

var (
	ErrEmptyName     = errors.New("name cannot be empty")
	ErrTooLongName   = errors.New("name too long (max 100 characters)")
	ErrInvalidUserID = errors.New("invalid user id")
)

// Category представляет сущность категоризации транзакций с опциональными родительско-дочерними отношениями.
// Категории принадлежат пользователям и используются для организации доходных и расходных транзакций.
type Category struct {
	ID        shared.ID
	Name      string
	OwnerID   shared.ID
	ParentID  shared.ID
	CreatedAt time.Time
}

// NewCategory создает новый экземпляр Category с валидированными параметрами.
// Возвращает ошибку если валидация имени не проходит или user ID невалиден.
func NewCategory(name string, uID, pID shared.ID) (Category, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return Category{}, ErrEmptyName
	}

	if len(name) > 100 {
		return Category{}, fmt.Errorf("%w: %s", ErrTooLongName, name)
	}

	if uID.IsZero() {
		return Category{}, fmt.Errorf("%w: %s", ErrInvalidUserID, uID)
	}

	c := Category{
		ID:        shared.NewID(),
		Name:      name,
		OwnerID:   uID,
		CreatedAt: time.Now(),
	}

	if !pID.IsZero() {
		c.ParentID = pID
	}

	return c, nil
}
