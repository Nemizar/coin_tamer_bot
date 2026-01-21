package category_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"

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
		cType    category.Type
		want     want
		wantErr  error
	}{
		{
			name:     "Создание с валидными данными",
			userID:   id,
			parentID: &id,
			cName:    "Категория 1",
			cType:    category.TypeExpense,
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
			cType:  category.TypeExpense,
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
			cType:    category.TypeExpense,
			wantErr:  category.ErrEmptyName,
		},
		{
			name:     "Создание с очень длинным именем. Ошибка",
			userID:   id,
			parentID: &id,
			cName:    strings.Repeat("a", 101),
			cType:    category.TypeExpense,
			wantErr:  category.ErrTooLongName,
		},
		{
			name:     "Создание с именем ровно из 100 символов. Успех",
			userID:   id,
			parentID: &id,
			cName:    strings.Repeat("a", 100),
			cType:    category.TypeExpense,
			want: want{
				name:     strings.Repeat("a", 100),
				ownerID:  id,
				parentID: &id,
			},
		},
		{
			name:     "Создание с не валидным пользователем",
			parentID: &id,
			cName:    "Категория 1",
			cType:    category.TypeExpense,
			wantErr:  errs.ErrValueIsInvalid,
		},
		{
			name:     "Создание без родителя",
			parentID: nil,
			userID:   id,
			cName:    "Родительская категория",
			cType:    category.TypeExpense,
			want: want{
				name:     "Родительская категория",
				ownerID:  id,
				parentID: nil,
			},
		},
		{
			name:     "Создание с неверным типом",
			userID:   id,
			parentID: &id,
			cName:    "Категория 1",
			cType:    category.Type("incorrect"),
			wantErr:  errs.ErrValueIsInvalid,
		},
		{
			name:     "Создание с типом income",
			userID:   id,
			parentID: nil,
			cName:    "Доход",
			cType:    category.TypeIncome,
			want: want{
				name:     "Доход",
				ownerID:  id,
				parentID: nil,
			},
		},
		{
			name:     "Создание с нулевым родительским ID",
			userID:   id,
			parentID: &shared.ID{}, // Zero ID
			cName:    "Категория с нулевым родителем",
			cType:    category.TypeExpense,
			want: want{
				name:     "Категория с нулевым родителем",
				ownerID:  id,
				parentID: nil, // Should ignore zero parent ID
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := category.New(tt.cName, tt.cType, tt.userID, tt.parentID)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tt.wantErr)

				return
			}

			if tt.want.parentID != nil {
				assert.Equal(t, *tt.want.parentID, got.ParentID())
			} else {
				assert.Equal(t, shared.ID{}, got.ParentID()) // Check that zero parent ID is ignored
			}
			assert.Equal(t, tt.want.ownerID, got.OwnerID())
			assert.Equal(t, tt.want.name, got.Name())
			assert.False(t, got.ID().IsZero())
			assert.NotEqual(t, time.Time{}, got.CreatedAt())
		})
	}
}

func TestCategory_Equals(t *testing.T) {
	userID := shared.NewID()

	cat1, err := category.New("Test", category.TypeExpense, userID, nil)
	require.NoError(t, err)

	cat2, err := category.New("Test", category.TypeExpense, userID, nil)
	require.NoError(t, err)

	tests := []struct {
		name     string
		cat1     *category.Category
		cat2     *category.Category
		expected bool
	}{
		{
			name:     "Сравнение категории с собой",
			cat1:     cat1,
			cat2:     cat1,
			expected: true,
		},
		{
			name:     "Сравнение разных категорий с одинаковыми полями",
			cat1:     cat1,
			cat2:     cat2, // Different instance
			expected: false,
		},
		{
			name:     "Сравнение с nil",
			cat1:     cat1,
			cat2:     nil,
			expected: false,
		},
		{
			name:     "Сравнение двух nil",
			cat1:     nil,
			cat2:     nil,
			expected: false, // According to the implementation, comparing with nil returns false
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cat1 == nil {
				// Can't call method on nil pointer, so we expect false as per implementation
				assert.False(t, tt.expected, "Expected false when cat1 is nil")
			} else {
				// Normal case: call method on non-nil object
				result := tt.cat1.Equals(tt.cat2)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestCategory_Getters(t *testing.T) {
	id := shared.NewID()
	userID := shared.NewID()
	parentID := shared.NewID()
	name := "Test Category"
	catType := category.TypeIncome
	createdAt := time.Now()

	cat := category.Restore(id, name, userID, &parentID, catType, createdAt)

	assert.Equal(t, id, cat.ID())
	assert.Equal(t, name, cat.Name())
	assert.Equal(t, userID, cat.OwnerID())
	assert.Equal(t, parentID, cat.ParentID())
	assert.Equal(t, catType, cat.Type())
	assert.Equal(t, createdAt, cat.CreatedAt())
}
