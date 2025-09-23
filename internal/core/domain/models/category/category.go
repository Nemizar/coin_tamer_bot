// Package category определяет доменную модель Category для организации транзакций.
// Обрабатывает создание категорий, валидацию и иерархические отношения.
package category

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/ddd"
)

var (
	ErrEmptyName     = errors.New("name cannot be empty")
	ErrTooLongName   = errors.New("name too long (max 100 characters)")
	ErrInvalidUserID = errors.New("invalid user id")
)

// Category представляет сущность категоризации транзакций с опциональными родительско-дочерними отношениями.
// Категории принадлежат пользователям и используются для организации доходных и расходных транзакций.
type Category struct {
	baseAggregate *ddd.BaseAggregate[shared.ID]
	name          string
	ownerID       shared.ID
	parentID      shared.ID
	createdAt     time.Time
}

// NewCategory создает новый экземпляр Category с валидированными параметрами.
// Возвращает ошибку если валидация имени не проходит или user id невалиден.
func NewCategory(name string, uID, pID shared.ID) (*Category, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, ErrEmptyName
	}

	if len(name) > 100 {
		return nil, fmt.Errorf("%w: %s", ErrTooLongName, name)
	}

	if uID.IsZero() {
		return nil, fmt.Errorf("%w: %s", ErrInvalidUserID, uID)
	}

	c := Category{
		baseAggregate: ddd.NewBaseAggregate(shared.NewID()),
		name:          name,
		ownerID:       uID,
		createdAt:     time.Now(),
	}

	if !pID.IsZero() {
		c.parentID = pID
	}

	return &c, nil
}

// ID возвращает уникальный идентификатор категории.
func (c Category) ID() shared.ID {
	return c.baseAggregate.ID()
}

// Name возвращает название категории.
func (c Category) Name() string {
	return c.name
}

// OwnerID возвращает идентификатор владельца категории.
func (c Category) OwnerID() shared.ID {
	return c.ownerID
}

// ParentID возвращает идентификатор родительской категории.
func (c Category) ParentID() shared.ID {
	return c.parentID
}

// CreatedAt возвращает временную метку создания категории.
func (c Category) CreatedAt() time.Time {
	return c.createdAt
}

// Equals сравнивает текущую категорию с другой категорией.
func (c Category) Equals(other *Category) bool {
	if other == nil {
		return false
	}

	return c.baseAggregate.Equal(other.baseAggregate)
}
