package user_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/user"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name      string
		username  string
		chatID    string
		wantError error
	}{
		{
			name:      "Валидный пользователь",
			username:  "Alice",
			wantError: nil,
			chatID:    "1",
		},
		{
			name:      "Невалидное имя - пустое",
			username:  "",
			chatID:    "1",
			wantError: errs.ErrValueIsRequired,
		},
		{
			name:      "Невалидное имя - состоит из пробелов (пустое)",
			username:  "   ",
			chatID:    "1",
			wantError: errs.ErrValueIsRequired,
		},
		{
			name:      "Невалидный идентификатор чата",
			username:  "Alice",
			chatID:    "",
			wantError: errs.ErrValueIsRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := user.New(tt.username, tt.chatID, user.ProviderTelegram)

			if tt.wantError != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantError), "expected %v, got %v", tt.wantError, err)
				assert.Nil(t, u)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.username, u.Name())
			assert.NotEqual(t, shared.ID{}, u.ID())
			assert.WithinDuration(t, time.Now(), u.CreatedAt(), time.Second)
			assert.NotEmpty(t, u.GetExternalIdentities())
			assert.Equal(t, "1", u.GetExternalIdentities()[0].ExternalID())
			assert.Equal(t, user.ProviderTelegram, u.GetExternalIdentities()[0].Provider())
			assert.Equal(t, u.ID(), u.GetExternalIdentities()[0].UserID())
		})
	}
}
