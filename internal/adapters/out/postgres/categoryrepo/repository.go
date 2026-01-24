package categoryrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

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
	stmt := `INSERT INTO categories (id, name, type, owner_id, parent_category_id, created_at)
			 VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := c.tracker.Tx().ExecContext(ctx, stmt, cat.ID(), cat.Name(), cat.Type(), cat.OwnerID(), cat.ParentID(), cat.CreatedAt())
	if err != nil {
		return fmt.Errorf("category repo create: %w", err)
	}

	return nil
}

func (c CategoryRepository) GetIncomeByUserID(ctx context.Context, userID shared.ID) ([]*category.Category, error) {
	stmt := `SELECT * FROM categories WHERE owner_id = $1 AND type = $2 AND parent_category_id = $3`

	return c.getByUserIDAndType(ctx, stmt, userID, category.TypeIncome)
}

func (c CategoryRepository) GetExpenseByUserID(ctx context.Context, userID shared.ID) ([]*category.Category, error) {
	stmt := `SELECT * FROM categories WHERE owner_id = $1 AND type = $2 AND parent_category_id != $3`

	return c.getByUserIDAndType(ctx, stmt, userID, category.TypeExpense)
}

func (c CategoryRepository) getByUserIDAndType(ctx context.Context, stmt string, userID shared.ID, t category.Type) ([]*category.Category, error) {
	rows, err := c.tracker.DB().QueryContext(ctx, stmt, userID, t, uuid.Nil)
	if err != nil {
		return nil, fmt.Errorf("category repo get by user id and type: %w", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			c.tracker.Logger().Error("category repo get by user id and type", "err", err.Error())
		}
	}(rows)

	var categories []*category.Category
	for rows.Next() {
		var model Model

		if err := rows.Scan(&model.ID, &model.Name, &model.OwnerID, &model.ParentID, &model.CategoryType, &model.CreatedAt); err != nil {
			return nil, fmt.Errorf("category repo get by user id and type: %w", err)
		}

		parentID := shared.ID{}
		if model.ParentID != uuid.Nil {
			parentID = shared.RestoreID(model.ParentID)
		}

		categories = append(categories, category.Restore(shared.RestoreID(model.ID), model.Name, shared.RestoreID(model.OwnerID), &parentID, model.CategoryType, model.CreatedAt))
	}

	return categories, nil
}

func (c CategoryRepository) HasCategoriesByUserID(ctx context.Context, userID shared.ID) (bool, error) {
	stmt := `SELECT EXISTS(SELECT 1 FROM categories WHERE owner_id = $1)`

	var exists bool
	err := c.tracker.DB().QueryRowContext(ctx, stmt, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("category repo has categories by user id: %w", err)
	}

	return exists, nil
}
