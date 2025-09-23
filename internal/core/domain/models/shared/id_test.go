package shared_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
)

func TestNewID(t *testing.T) {
	id := shared.NewID()

	require.False(t, id.IsZero(), "NewID() должен возвращать непустой id")
	require.NotEqual(t, uuid.Nil, id.Value())
	require.Equal(t, id.String(), id.Value().String())
}

func TestNewIDFromString_Valid(t *testing.T) {
	u := uuid.New()
	id, err := shared.NewIDFromString(u.String())

	require.NoError(t, err)
	require.Equal(t, u, id.Value())
	require.Equal(t, u.String(), id.String())
}

func TestNewIDFromString_Invalid(t *testing.T) {
	_, err := shared.NewIDFromString("not-a-uuid")

	require.Error(t, err, "ожидается ошибка при парсинге некорректного UUID")
}

func TestIsZero(t *testing.T) {
	var zero shared.ID
	require.True(t, zero.IsZero(), "нулевой id должен быть пустым")

	id := shared.NewID()
	require.False(t, id.IsZero(), "сгенерированный id не должен быть пустым")
}

func TestValue(t *testing.T) {
	id := shared.NewID()
	val := id.Value()

	require.IsType(t, uuid.UUID{}, val)
	require.Equal(t, id.String(), val.String())
}
