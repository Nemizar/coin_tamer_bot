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
		provider  user.Provider
		wantError error
	}{
		{
			name:      "Валидный пользователь",
			username:  "Alice",
			chatID:    "1",
			provider:  user.ProviderTelegram,
			wantError: nil,
		},
		{
			name:      "Невалидное имя - пустое",
			username:  "",
			chatID:    "1",
			provider:  user.ProviderTelegram,
			wantError: errs.ErrValueIsRequired,
		},
		{
			name:      "Невалидное имя - состоит из пробелов (пустое)",
			username:  "   ",
			chatID:    "1",
			provider:  user.ProviderTelegram,
			wantError: errs.ErrValueIsRequired,
		},
		{
			name:      "Невалидный идентификатор чата",
			username:  "Alice",
			chatID:    "",
			provider:  user.ProviderTelegram,
			wantError: errs.ErrValueIsRequired,
		},
		{
			name:      "Валидный пользователь с другим провайдером",
			username:  "Bob",
			chatID:    "987654321",
			provider:  user.ProviderTelegram,
			wantError: nil,
		},
		{
			name:      "Имя пользователя со специальными символами",
			username:  "user_name-123",
			chatID:    "123456789",
			provider:  user.ProviderTelegram,
			wantError: nil,
		},
		{
			name:      "Имя пользователя с юникод символами",
			username:  "Юзер",
			chatID:    "123456789",
			provider:  user.ProviderTelegram,
			wantError: nil,
		},
		{
			name:      "ID чата с ведущими нулями",
			username:  "TestUser",
			chatID:    "000123456",
			provider:  user.ProviderTelegram,
			wantError: nil,
		},
		{
			name:      "ID чата в виде строки чисел",
			username:  "TestUser",
			chatID:    "999888777666",
			provider:  user.ProviderTelegram,
			wantError: nil,
		},
		{
			name:      "Имя пользователя с завершающими пробелами (должно сохраняться)",
			username:  "UserWithSpaces   ",
			chatID:    "123456789",
			provider:  user.ProviderTelegram,
			wantError: nil,
		},
		{
			name:      "ID чата со значением '0' (должен завершиться ошибкой)",
			username:  "TestUser",
			chatID:    "0",
			provider:  user.ProviderTelegram,
			wantError: errs.ErrValueIsRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := user.New(tt.username, tt.chatID, tt.provider)

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
			assert.NotEmpty(t, u.GetExternalIdentity())
			assert.Equal(t, tt.chatID, u.GetExternalIdentity().ExternalID())
			assert.Equal(t, tt.provider, u.GetExternalIdentity().Provider())
			assert.Equal(t, u.ID(), u.GetExternalIdentity().UserID())
		})
	}
}

func TestUser_AddExternalIdentity(t *testing.T) {
	u, err := user.New("TestUser", "123456789", user.ProviderTelegram)
	require.NoError(t, err)

	// Test adding a duplicate external identity
	err = u.AddExternalIdentity(u.GetExternalIdentity())
	assert.Error(t, err)
	// The actual error is ErrValueIsInvalid for the externalIdentity field
	assert.Contains(t, err.Error(), "invalid")

	// Test adding nil external identity
	err = u.AddExternalIdentity(nil)
	assert.Error(t, err)
	// The actual error is ErrValueIsRequired for externalIdentity field
	assert.Contains(t, err.Error(), "invalid") // The error message is "value is invalid: externalIdentity"

	// Test adding external identity with wrong user ID
	wrongUserID := shared.NewID()
	wrongEI, err := user.NewExternalIdentity(wrongUserID, user.ProviderTelegram, "987654321")
	require.NoError(t, err)

	err = u.AddExternalIdentity(wrongEI)
	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrValueIsInvalid)
}

func TestUser_Equals(t *testing.T) {
	u1, err := user.New("User1", "123456789", user.ProviderTelegram)
	require.NoError(t, err)

	u2, err := user.New("User2", "987654321", user.ProviderTelegram)
	require.NoError(t, err)

	tests := []struct {
		name     string
		u1       *user.User
		u2       *user.User
		expected bool
	}{
		{
			name:     "Один и тот же экземпляр пользователя",
			u1:       u1,
			u2:       u1,
			expected: true,
		},
		{
			name:     "Разные пользователи",
			u1:       u1,
			u2:       u2,
			expected: false,
		},
		{
			name:     "Один пользователь равен nil",
			u1:       u1,
			u2:       nil,
			expected: false,
		},
		{
			name:     "Оба пользователя равны nil",
			u1:       nil,
			u2:       nil,
			expected: false, // According to the implementation, comparing with nil returns false
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.u1.Equals(tt.u2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUser_Restore(t *testing.T) {
	id := shared.NewID()
	name := "Restored User"
	createdAt := time.Now()

	u := user.Restore(id, name, createdAt)

	assert.Equal(t, id, u.ID())
	assert.Equal(t, name, u.Name())
	assert.Equal(t, createdAt, u.CreatedAt())
	assert.Nil(t, u.GetExternalIdentity()) // Restored user should have no external identity initially
}
