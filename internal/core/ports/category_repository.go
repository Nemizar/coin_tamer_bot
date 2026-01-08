package ports

import (
	"context"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/category"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *category.Category) error
	GetIncomeByUserID(ctx context.Context, userID shared.ID) ([]*category.Category, error)
	GetExpenseByUserID(ctx context.Context, userID shared.ID) ([]*category.Category, error)
}
