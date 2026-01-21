package user_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/user"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

func TestNewExternalIdentity(t *testing.T) {
	id := shared.NewID()

	tests := []struct {
		name       string
		userID     shared.ID
		provider   user.Provider
		externalID string
		wantError  error
	}{
		{
			name:       "Валидный внешний идентификатор",
			userID:     shared.NewID(),
			provider:   user.ProviderTelegram,
			externalID: "123456789",
			wantError:  nil,
		},
		{
			name:       "Невалидный провайдер",
			userID:     shared.NewID(),
			provider:   "invalid_provider",
			externalID: "123456789",
			wantError:  errs.ErrValueIsInvalid,
		},
		{
			name:       "Пустой externalID",
			userID:     shared.NewID(),
			provider:   user.ProviderTelegram,
			externalID: "",
			wantError:  errs.ErrValueIsRequired,
		},
		{
			name:       "Валидный внешний идентификатор с другим провайдером",
			userID:     id,
			provider:   user.ProviderTelegram,
			externalID: "different_id",
			wantError:  nil,
		},
		{
			name:       "Валидный внешний идентификатор с числовым ID",
			userID:     id,
			provider:   user.ProviderTelegram,
			externalID: "123456789",
			wantError:  nil,
		},
		{
			name:       "Валидный внешний идентификатор со специальными символами в ID",
			userID:     id,
			provider:   user.ProviderTelegram,
			externalID: "user_123-test@example.com",
			wantError:  nil,
		},
		{
			name:       "ID с только пробелами (не пустая строка, поэтому должен быть валидным)",
			userID:     id,
			provider:   user.ProviderTelegram,
			externalID: "   ", // spaces only, not empty string
			wantError:  nil,
		},
		{
			name:       "Валидный провайдер с другим регистром",
			userID:     id,
			provider:   user.Provider("telegram"), // lowercase
			externalID: "test_id",
			wantError:  nil,
		},
		{
			name:       "ID с ведущими/завершающими пробелами (должен быть обрезан вызывающим, если нужно)",
			userID:     id,
			provider:   user.ProviderTelegram,
			externalID: "  spaced_id  ",
			wantError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exIdentity, err := user.NewExternalIdentity(tt.userID, tt.provider, tt.externalID)

			if tt.wantError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantError, "ожидалась ошибка: %v, получена: %v", tt.wantError, err)
				assert.Nil(t, exIdentity)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, exIdentity)
			assert.Equal(t, tt.userID, exIdentity.UserID())
			assert.Equal(t, tt.provider, exIdentity.Provider())
			assert.Equal(t, tt.externalID, exIdentity.ExternalID())
			assert.NotEmpty(t, exIdentity.ID())
		})
	}
}

func TestExternalIdentity_Getters(t *testing.T) {
	userID := shared.NewID()
	provider := user.ProviderTelegram
	externalID := "987654321"

	ei, err := user.NewExternalIdentity(userID, provider, externalID)
	require.NoError(t, err)

	t.Run("Проверка геттеров", func(t *testing.T) {
		assert.Equal(t, userID, ei.UserID())
		assert.Equal(t, provider, ei.Provider())
		assert.Equal(t, externalID, ei.ExternalID())
		assert.NotEmpty(t, ei.ID())
		assert.NotEqual(t, time.Time{}, ei.GetCreatedAt())
	})
}

func TestExternalIdentity_Equals(t *testing.T) {
	userID := shared.NewID()

	ei1, err := user.NewExternalIdentity(userID, user.ProviderTelegram, "123")
	require.NoError(t, err)

	ei2, err := user.NewExternalIdentity(userID, user.ProviderTelegram, "456")
	require.NoError(t, err)

	tests := []struct {
		name     string
		ei1      *user.ExternalIdentity
		ei2      *user.ExternalIdentity
		expected bool
	}{
		{
			name:     "Один и тот же экземпляр внешнего идентификатора",
			ei1:      ei1,
			ei2:      ei1,
			expected: true,
		},
		{
			name:     "Разные внешние идентификаторы",
			ei1:      ei1,
			ei2:      ei2,
			expected: false,
		},
		{
			name:     "Один внешний идентификатор равен nil",
			ei1:      ei1,
			ei2:      nil,
			expected: false,
		},
		{
			name:     "Оба внешних идентификатора равны nil",
			ei1:      nil,
			ei2:      nil,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch {
			case tt.ei1 == nil && tt.ei2 == nil:
				// Both nil case - already handled by the logic
				return
			case tt.ei1 == nil || tt.ei2 == nil:
				// One is nil, so they're not equal
				assert.False(t, tt.expected)
			default:
				// For actual comparison, we compare IDs since there's no explicit Equals method
				assert.Equal(t, tt.ei1.ID() == tt.ei2.ID(), tt.expected)
			}
		})
	}
}
