package queries

import (
	"context"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/category"
	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

type GetUserCategoriesByTypeQueryHandler interface {
	Handle(ctx context.Context, query GetUserCategoriesByTypeQuery) ([]*category.Category, error)
}

type getUserCategoriesByTypeQueryHandler struct {
	uow ports.UnitOfWork
}

func NewGetUserCategoriesByTypeQueryHandler(uow ports.UnitOfWork) (GetUserCategoriesByTypeQueryHandler, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequiredError("uow")
	}

	return &getUserCategoriesByTypeQueryHandler{uow: uow}, nil
}

func (h getUserCategoriesByTypeQueryHandler) Handle(ctx context.Context, query GetUserCategoriesByTypeQuery) ([]*category.Category, error) {
	var (
		result []*category.Category
		err    error
	)

	if query.CategoryType() == category.TypeIncome {
		result, err = h.uow.CategoryRepository().GetIncomeByUserID(ctx, query.UserID())
		if err != nil {
			return nil, err
		}
	} else {
		result, err = h.uow.CategoryRepository().GetExpenseByUserID(ctx, query.UserID())
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
