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
		wantError error
	}{
		{
			name:      "Валидный пользователь",
			username:  "Alice",
			wantError: nil,
		},
		{
			name:      "Невалидное имя - пустое",
			username:  "",
			wantError: errs.ErrValueIsRequired,
		},
		{
			name:      "Невалидное имя - состоит из пробелов (пустое)",
			username:  "   ",
			wantError: errs.ErrValueIsRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := user.New(tt.username)

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
		})
	}
}
