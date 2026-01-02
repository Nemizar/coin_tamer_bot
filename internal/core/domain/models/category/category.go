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

type Category struct {
	baseAggregate *ddd.BaseAggregate[shared.ID]
	name          string
	ownerID       shared.ID
	parentID      shared.ID
	createdAt     time.Time
}

func New(name string, uID shared.ID, pID *shared.ID) (*Category, error) {
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

	if pID != nil && !pID.IsZero() {
		c.parentID = *pID
	}

	return &c, nil
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

func (c Category) CreatedAt() time.Time {
	return c.createdAt
}

func (c Category) Equals(other *Category) bool {
	if other == nil {
		return false
	}

	return c.baseAggregate.Equal(other.baseAggregate)
}
