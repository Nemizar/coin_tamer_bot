package queries

import (
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/category"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
)

type GetUserCategoriesByTypeQuery interface {
	UserID() shared.ID
	CategoryType() category.Type
}

type getUserCategoriesByTypeQuery struct {
	userID       shared.ID
	categoryType category.Type
}

func NewGetUserCategoriesByType(userID shared.ID, categoryType category.Type) GetUserCategoriesByTypeQuery {
	return &getUserCategoriesByTypeQuery{
		userID:       userID,
		categoryType: categoryType,
	}
}

func (g getUserCategoriesByTypeQuery) UserID() shared.ID {
	return g.userID
}

func (g getUserCategoriesByTypeQuery) CategoryType() category.Type {
	return g.categoryType
}
