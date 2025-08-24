package user_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Nemizar/coin_tamer_bot/internal/domain/shared"
	"github.com/Nemizar/coin_tamer_bot/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name      string
		chatID    int64
		username  string
		wantError error
	}{
		{
			name:      "Валидный пользователь",
			chatID:    12345,
			username:  "Alice",
			wantError: nil,
		},
		{
			name:      "Невалидный идентификатор чата",
			chatID:    0,
			username:  "Alice",
			wantError: user.ErrInvalidChatID,
		},
		{
			name:      "Невалидное имя - пустое",
			chatID:    12345,
			username:  "",
			wantError: user.ErrInvalidName,
		},
		{
			name:      "Невалидное имя - состоит из пробелов (пустое)",
			chatID:    12345,
			username:  "   ",
			wantError: user.ErrInvalidName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := user.NewUser(tt.chatID, tt.username)

			if tt.wantError != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantError), "expected %v, got %v", tt.wantError, err)
				assert.Equal(t, user.User{}, u)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.chatID, u.ChatID)
			assert.Equal(t, tt.username, u.Name)
			assert.NotEqual(t, shared.ID{}, u.ID)
			assert.WithinDuration(t, time.Now(), u.CreatedAT, time.Second)
		})
	}
}
