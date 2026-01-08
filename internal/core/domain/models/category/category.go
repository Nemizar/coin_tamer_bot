package category

import (
	"errors"
	"strings"
	"time"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/ddd"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

var (
	ErrEmptyName   = errors.New("name cannot be empty")
	ErrTooLongName = errors.New("name too long (max 100 characters)")
)

type Category struct {
	baseAggregate *ddd.BaseAggregate[shared.ID]
	name          string
	ownerID       shared.ID
	parentID      shared.ID
	categoryType  Type
	createdAt     time.Time
}

func New(name string, categoryType Type, uID shared.ID, pID *shared.ID) (*Category, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, ErrEmptyName
	}

	if len(name) > 100 {
		return nil, errs.NewValueIsInvalidErrorWithCause("name", ErrTooLongName)
	}

	if uID.IsZero() {
		return nil, errs.NewValueIsInvalidError("ownerID")
	}

	if !categoryType.IsValid() {
		return nil, errs.NewValueIsInvalidError("categoryType")
	}

	c := Category{
		baseAggregate: ddd.NewBaseAggregate(shared.NewID()),
		name:          name,
		ownerID:       uID,
		createdAt:     time.Now(),
		categoryType:  categoryType,
	}

	if pID != nil && !pID.IsZero() {
		c.parentID = *pID
	}

	return &c, nil
}

func Restore(id shared.ID, name string, ownerID shared.ID, parentID *shared.ID, categoryType Type, createdAt time.Time) *Category {
	c := Category{
		baseAggregate: ddd.NewBaseAggregate(id),
		name:          name,
		ownerID:       ownerID,
		categoryType:  categoryType,
		createdAt:     createdAt,
	}

	if parentID != nil && !parentID.IsZero() {
		c.parentID = *parentID
	}

	return &c
}

func (c Category) ID() shared.ID {
	return c.baseAggregate.ID()
}

func (c Category) Name() string {
	return c.name
}

func (c Category) OwnerID() shared.ID {
	return c.ownerID
}

func (c Category) ParentID() shared.ID {
	return c.parentID
}

func (c Category) Type() Type {
	return c.categoryType
}

func (c Category) CreatedAt() time.Time {
	return c.createdAt
}

func (c Category) Equals(other *Category) bool {
	if other == nil {
		return false
	}

	return c.baseAggregate.Equal(other.baseAggregate)
}
