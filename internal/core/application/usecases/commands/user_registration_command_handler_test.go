package commands_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Nemizar/coin_tamer_bot/mocks/core/portsmocks"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/user"

	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"

	"github.com/Nemizar/coin_tamer_bot/internal/core/application/usecases/commands"
)

func TestUserRegistrationCommandHandler_Success(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	u, err := user.New("test", "123", user.ProviderTelegram)
	require.Nil(t, err)

	cmd, err := commands.NewUserRegistrationCommand(u.Name(), u.GetExternalIdentity().ExternalID(), u.GetExternalIdentity().Provider())
	assert.Nil(t, err)

	var captureObj *user.User
	userRepoMock := &portsmocks.UserRepositoryMock{}
	userRepoMock.
		EXPECT().
		FindByExternalProvider(ctx, u.GetExternalIdentity().Provider(), u.GetExternalIdentity().ExternalID()).
		Return(nil, nil).
		Once()
	userRepoMock.
		EXPECT().
		Create(ctx, mock.AnythingOfType("*user.User")).
		Run(func(ctx context.Context, u *user.User) {
			captureObj = u
		}).
		Return(nil).
		Once()

	uowMock := &portsmocks.UnitOfWorkMock{}
	uowMock.
		On("UserRepository").
		Return(userRepoMock)
	uowMock.
		EXPECT().
		Begin(ctx).
		Return(nil)
	uowMock.
		EXPECT().
		RollbackUnlessCommitted().
		Return(nil)
	uowMock.
		EXPECT().
		Commit(ctx).
		Return(nil)

	handler, err := commands.NewUserRegistrationCommandHandler(logger, uowMock)
	require.Nil(t, err)

	err = handler.Handle(ctx, cmd)
	assert.Nil(t, err)
	assert.Equal(t, "test", captureObj.Name())
	assert.NotEqual(t, uuid.Nil, captureObj.ID())
	assert.Equal(t, user.ProviderTelegram, captureObj.GetExternalIdentity().Provider())
	assert.Equal(t, "123", captureObj.GetExternalIdentity().ExternalID())
}

func TestUserRegistrationCommandHandler_Failure_EmptyName(t *testing.T) {
	cmd, err := commands.NewUserRegistrationCommand("", "123", user.ProviderTelegram)
	assert.Nil(t, cmd)
	assert.ErrorIs(t, err, errs.ErrValueIsRequired)
}

func TestUserRegistrationCommandHandler_Failure_EmptyTelegramChatID(t *testing.T) {
	cmd, err := commands.NewUserRegistrationCommand("test", "0", user.ProviderTelegram)
	assert.Nil(t, cmd)
	assert.ErrorIs(t, err, errs.ErrValueIsRequired)
}

func TestUserRegistrationCommandHandler_Idempotent(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	u, err := user.New("test", "123", user.ProviderTelegram)
	require.NoError(t, err)

	cmd, err := commands.NewUserRegistrationCommand(
		u.Name(),
		u.GetExternalIdentity().ExternalID(),
		u.GetExternalIdentity().Provider(),
	)
	require.NoError(t, err)

	userRepoMock := &portsmocks.UserRepositoryMock{}
	userRepoMock.
		EXPECT().
		FindByExternalProvider(
			ctx,
			user.ProviderTelegram,
			"123",
		).
		Return(nil, nil).
		Once()
	userRepoMock.
		EXPECT().
		Create(ctx, mock.AnythingOfType("*user.User")).
		Return(nil).
		Once()
	userRepoMock.
		EXPECT().
		FindByExternalProvider(
			ctx,
			user.ProviderTelegram,
			"123",
		).
		Return(u, nil).
		Once()

	uowMock := &portsmocks.UnitOfWorkMock{}
	uowMock.
		On("UserRepository").
		Return(userRepoMock)
	uowMock.
		EXPECT().
		Begin(ctx).
		Return(nil).
		Times(2)
	uowMock.
		EXPECT().
		RollbackUnlessCommitted().
		Return(nil).
		Times(2)
	uowMock.
		EXPECT().
		Commit(ctx).
		Return(nil).
		Once()

	handler, err := commands.NewUserRegistrationCommandHandler(logger, uowMock)
	require.NoError(t, err)

	err = handler.Handle(ctx, cmd)
	require.NoError(t, err)

	err = handler.Handle(ctx, cmd)
	require.NoError(t, err)
}
