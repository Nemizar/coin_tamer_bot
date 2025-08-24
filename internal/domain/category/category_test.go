package category_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nemizar/coin_tamer_bot/internal/domain/category"
	"github.com/Nemizar/coin_tamer_bot/internal/domain/shared"
)

func TestNewCategory(t *testing.T) {
	id := shared.NewID()

	tests := []struct {
		name     string
		userID   shared.ID
		parentID shared.ID
		cName    string
		want     category.Category
		wantErr  error
	}{
		{
			name:     "Создание с валидными данными",
			userID:   id,
			parentID: id,
			cName:    "Категория 1",
			want: category.Category{
				Name:     "Категория 1",
				OwnerID:  id,
				ParentID: id,
			},
		},
		{
			name:   "Создание родительской категории",
			userID: id,
			cName:  "Категория 1",
			want: category.Category{
				Name:    "Категория 1",
				OwnerID: id,
			},
		},
		{
			name:     "Создание без измени. Ошибка",
			userID:   id,
			parentID: id,
			cName:    "",
			wantErr:  category.ErrEmptyName,
		},
		{
			name:     "Создание с очень длинным именем. Ошибка",
			userID:   id,
			parentID: id,
			cName:    strings.Repeat("a", 101),
			wantErr:  category.ErrTooLongName,
		},
		{
			name:     "Создание с не валидным пользователем",
			parentID: id,
			cName:    "Категория 1",
			wantErr:  category.ErrInvalidUserID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := category.NewCategory(tt.cName, tt.userID, tt.parentID)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)

				return
			}

			assert.Equal(t, tt.want.ParentID, got.ParentID)
			assert.Equal(t, tt.want.OwnerID, got.OwnerID)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.False(t, got.ID.IsZero())
			assert.NotEqual(t, time.Time{}, got.CreatedAt)
		})
	}
}
