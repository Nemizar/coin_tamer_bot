package categoryrepo

import (
	"context"
	"fmt"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/category"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

type CategoryRepository struct {
	tracker Tracker
}

func NewCategoryRepository(tracker Tracker) (ports.CategoryRepository, error) {
	if tracker == nil {
		return nil, errs.NewValueIsRequiredError("tracker")
	}

	return &CategoryRepository{tracker: tracker}, nil
}

func (c CategoryRepository) Create(ctx context.Context, cat *category.Category) error {
	stmt := `INSERT INTO categories (id, name, owner_id, parent_category_id, created_at)
			 VALUES ($1, $2, $3, $4, $5)`
	_, err := c.tracker.Tx().ExecContext(ctx, stmt, cat.ID(), cat.Name(), cat.OwnerID(), cat.ParentID(), cat.CreatedAt())
	if err != nil {
		return fmt.Errorf("category repo create: %w", err)
	}

	return nil
}

func (c CategoryRepository) GetAll(ctx context.Context) ([]*category.Category, error) {
	panic("implement me")
}

func (c CategoryRepository) Update(ctx context.Context, category *category.Category) error {
	panic("implement me")
}

func (c CategoryRepository) Delete(ctx context.Context, id shared.ID) error {
	panic("implement me")
}
