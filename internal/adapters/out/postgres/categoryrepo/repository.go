package categoryrepo

import (
	"context"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/category"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"
)

type CategoryRepository struct {
	uow Tracker
}

func NewCategoryRepository(uow Tracker) ports.CategoryRepository {
	return &CategoryRepository{uow: uow}
}

func (c CategoryRepository) Create(ctx context.Context, category *category.Category) error {
	panic("implement me")
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
