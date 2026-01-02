package category_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/category"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
)

func TestNewCategory(t *testing.T) {
	id := shared.NewID()
	type want struct {
		name     string
		ownerID  shared.ID
		parentID *shared.ID
	}

	tests := []struct {
		name     string
		userID   shared.ID
		parentID *shared.ID
		cName    string
		want     want
		wantErr  error
	}{
		{
			name:     "Создание с валидными данными",
			userID:   id,
			parentID: &id,
			cName:    "Категория 1",
			want: want{
				name:     "Категория 1",
				ownerID:  id,
				parentID: &id,
			},
		},
		{
			name:   "Создание родительской категории",
			userID: id,
			cName:  "Категория 1",
			want: want{
				name:    "Категория 1",
				ownerID: id,
			},
		},
		{
			name:     "Создание без измени. Ошибка",
			userID:   id,
			parentID: &id,
			cName:    "",
			wantErr:  category.ErrEmptyName,
		},
		{
			name:     "Создание с очень длинным именем. Ошибка",
			userID:   id,
			parentID: &id,
			cName:    strings.Repeat("a", 101),
			wantErr:  category.ErrTooLongName,
		},
		{
			name:     "Создание с не валидным пользователем",
			parentID: &id,
			cName:    "Категория 1",
			wantErr:  category.ErrInvalidUserID,
		},
		{
			name:     "Создание без родителя",
			parentID: nil,
			userID:   id,
			cName:    "Родительская категория",
			want: want{
				name:     "Родительская категория",
				ownerID:  id,
				parentID: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := category.New(tt.cName, tt.userID, tt.parentID)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)

				return
			}

			if tt.want.parentID != nil {
				assert.Equal(t, *tt.want.parentID, got.ParentID())
			}
			assert.Equal(t, tt.want.ownerID, got.OwnerID())
			assert.Equal(t, tt.want.name, got.Name())
			assert.False(t, got.ID().IsZero())
			assert.NotEqual(t, time.Time{}, got.CreatedAt())
		})
	}
}
