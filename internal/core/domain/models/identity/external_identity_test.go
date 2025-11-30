package identity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/identity"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

func TestNewExternalIdentity(t *testing.T) {
	tests := []struct {
		name       string
		userID     shared.ID
		provider   identity.Provider
		externalID string
		wantError  error
	}{
		{
			name:       "Валидный внешний идентификатор",
			userID:     shared.NewID(),
			provider:   identity.ProviderTelegram,
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
			provider:   identity.ProviderTelegram,
			externalID: "",
			wantError:  errs.ErrValueIsRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exIdentity, err := identity.NewExternalIdentity(tt.userID, tt.provider, tt.externalID)

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
	provider := identity.ProviderTelegram
	externalID := "987654321"

	ei, err := identity.NewExternalIdentity(userID, provider, externalID)
	require.NoError(t, err)

	t.Run("Проверка геттеров", func(t *testing.T) {
		assert.Equal(t, userID, ei.UserID())
		assert.Equal(t, provider, ei.Provider())
		assert.Equal(t, externalID, ei.ExternalID())
		assert.NotEmpty(t, ei.ID())
	})
}
