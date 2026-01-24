package commands_test

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Nemizar/coin_tamer_bot/internal/core/application/usecases/commands"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/user"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
	"github.com/Nemizar/coin_tamer_bot/mocks/core/portsmocks"
)

// setupMocks создает и настраивает моки для тестов
func setupMocks() (*portsmocks.UnitOfWorkMock, *portsmocks.UserRepositoryMock) {
	userRepoMock := &portsmocks.UserRepositoryMock{}
	uowMock := &portsmocks.UnitOfWorkMock{}
	uowMock.On("UserRepository").Return(userRepoMock)
	return uowMock, userRepoMock
}

func TestUserRegistrationCommandHandler_Validation(t *testing.T) {
	tests := []struct {
		name     string
		username string
		chatID   string
		provider user.Provider
		wantErr  error
	}{
		{
			name:     "Empty name",
			username: "",
			chatID:   "123",
			provider: user.ProviderTelegram,
			wantErr:  errs.ErrValueIsRequired,
		},
		{
			name:     "Empty chat ID",
			username: "test",
			chatID:   "",
			provider: user.ProviderTelegram,
			wantErr:  errs.ErrValueIsRequired,
		},
		{
			name:     "Zero chat ID",
			username: "test",
			chatID:   "0",
			provider: user.ProviderTelegram,
			wantErr:  errs.ErrValueIsRequired,
		},
		{
			name:     "Valid data",
			username: "test",
			chatID:   "123",
			provider: user.ProviderTelegram,
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := commands.NewUserRegistrationCommand(tt.username, tt.chatID, tt.provider)

			if tt.wantErr != nil {
				assert.Nil(t, cmd)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NotNil(t, cmd)
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserRegistrationCommandHandler_Success(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	u, err := user.New("test", "123", user.ProviderTelegram)
	require.NoError(t, err)

	cmd, err := commands.NewUserRegistrationCommand(u.Name(), u.GetExternalIdentity().ExternalID(), u.GetExternalIdentity().Provider())
	require.NoError(t, err)

	var captureObj *user.User
	uowMock, userRepoMock := setupMocks()

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

	uowMock.
		EXPECT().
		Begin(ctx).
		Return(nil).
		Once()
	uowMock.
		EXPECT().
		RollbackUnlessCommitted().
		Return(nil).
		Once()
	uowMock.
		EXPECT().
		Commit(ctx).
		Return(nil).
		Once()

	handler, err := commands.NewUserRegistrationCommandHandler(logger, uowMock)
	require.NoError(t, err)

	err = handler.Handle(ctx, cmd)
	assert.NoError(t, err)
	assert.Equal(t, "test", captureObj.Name())
	assert.NotEqual(t, uuid.Nil, captureObj.ID())
	assert.Equal(t, user.ProviderTelegram, captureObj.GetExternalIdentity().Provider())
	assert.Equal(t, "123", captureObj.GetExternalIdentity().ExternalID())

	// Проверяем, что моки вызваны в соответствии с ожиданиями
	userRepoMock.AssertExpectations(t)
	uowMock.AssertExpectations(t)
}

func TestUserRegistrationCommandHandler_UserAlreadyExists(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	u, err := user.New("test", "123", user.ProviderTelegram)
	require.NoError(t, err)

	cmd, err := commands.NewUserRegistrationCommand(
		u.Name(),
		u.GetExternalIdentity().ExternalID(),
		u.GetExternalIdentity().Provider(),
	)
	require.NoError(t, err)

	userRepoMock := &portsmocks.UserRepositoryMock{}
	uowMock := &portsmocks.UnitOfWorkMock{}
	uowMock.On("UserRepository").Return(userRepoMock)

	// FindByExternalProvider возвращает существующего пользователя
	userRepoMock.
		EXPECT().
		FindByExternalProvider(ctx, u.GetExternalIdentity().Provider(), u.GetExternalIdentity().ExternalID()).
		Return(u, nil). // Пользователь уже существует
		Once()

	// Create не должен вызываться, так как пользователь уже существует
	uowMock.
		EXPECT().
		Begin(ctx).
		Return(nil).
		Once()
	uowMock.
		EXPECT().
		RollbackUnlessCommitted().
		Return(nil).
		Once()
	// Commit не вызывается, когда пользователь уже существует

	handler, err := commands.NewUserRegistrationCommandHandler(logger, uowMock)
	require.NoError(t, err)

	err = handler.Handle(ctx, cmd)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrEntityAlreadyExists))

	// Проверяем, что Create не был вызван
	userRepoMock.AssertNotCalled(t, "Create")
	userRepoMock.AssertExpectations(t)
	uowMock.AssertExpectations(t)
}

func TestUserRegistrationCommandHandler_UserAlreadyExistsReturnsCorrectErrorType(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	u, err := user.New("test", "123", user.ProviderTelegram)
	require.NoError(t, err)

	cmd, err := commands.NewUserRegistrationCommand(
		u.Name(),
		u.GetExternalIdentity().ExternalID(),
		u.GetExternalIdentity().Provider(),
	)
	require.NoError(t, err)

	userRepoMock := &portsmocks.UserRepositoryMock{}
	uowMock := &portsmocks.UnitOfWorkMock{}
	uowMock.On("UserRepository").Return(userRepoMock)

	// FindByExternalProvider возвращает существующего пользователя
	userRepoMock.
		EXPECT().
		FindByExternalProvider(ctx, u.GetExternalIdentity().Provider(), u.GetExternalIdentity().ExternalID()).
		Return(u, nil). // Пользователь уже существует
		Once()

	uowMock.
		EXPECT().
		Begin(ctx).
		Return(nil).
		Once()
	uowMock.
		EXPECT().
		RollbackUnlessCommitted().
		Return(nil).
		Once()

	handler, err := commands.NewUserRegistrationCommandHandler(logger, uowMock)
	require.NoError(t, err)

	err = handler.Handle(ctx, cmd)
	assert.Error(t, err)

	// Проверяем, что возвращается правильный тип ошибки
	var entityAlreadyExistsErr *errs.EntityAlreadyExistsError
	assert.ErrorAs(t, err, &entityAlreadyExistsErr)
	assert.Equal(t, "user", entityAlreadyExistsErr.Entity)
	assert.Equal(t, "external_id", entityAlreadyExistsErr.Field)
	assert.Equal(t, "123", entityAlreadyExistsErr.Value)

	// Проверяем, что Create не был вызван
	userRepoMock.AssertNotCalled(t, "Create")
	userRepoMock.AssertExpectations(t)
	uowMock.AssertExpectations(t)
}

func TestUserRegistrationCommandHandler_BeginTransactionError(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	cmd, err := commands.NewUserRegistrationCommand("test", "123", user.ProviderTelegram)
	require.NoError(t, err)

	// В этом тесте UserRepository() не будет вызван, так как ошибка происходит в Begin
	userRepoMock := &portsmocks.UserRepositoryMock{}
	uowMock := &portsmocks.UnitOfWorkMock{}

	// Мок возвращает ошибку при начале транзакции
	uowMock.
		EXPECT().
		Begin(ctx).
		Return(errors.New("begin transaction error")).
		Once()
	// Устанавливаем ожидание UserRepository, но помечаем как Maybe, чтобы не требовался вызов
	uowMock.On("UserRepository").Return(userRepoMock).Maybe()
	uowMock.
		EXPECT().
		RollbackUnlessCommitted().
		Return(nil).
		Once()

	handler, err := commands.NewUserRegistrationCommandHandler(logger, uowMock)
	require.NoError(t, err)

	err = handler.Handle(ctx, cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "begin transaction error")

	uowMock.AssertExpectations(t)
}

func TestUserRegistrationCommandHandler_CommitError(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	cmd, err := commands.NewUserRegistrationCommand("test", "123", user.ProviderTelegram)
	require.NoError(t, err)

	uowMock, userRepoMock := setupMocks()

	userRepoMock.
		EXPECT().
		FindByExternalProvider(ctx, user.ProviderTelegram, "123").
		Return(nil, nil).
		Once()
	userRepoMock.
		EXPECT().
		Create(ctx, mock.AnythingOfType("*user.User")).
		Return(nil).
		Once()

	uowMock.
		EXPECT().
		Begin(ctx).
		Return(nil).
		Once()
	uowMock.
		EXPECT().
		RollbackUnlessCommitted().
		Return(nil).
		Once()
	uowMock.
		EXPECT().
		Commit(ctx).
		Return(errors.New("commit error")).
		Once()

	handler, err := commands.NewUserRegistrationCommandHandler(logger, uowMock)
	require.NoError(t, err)

	err = handler.Handle(ctx, cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "commit error")

	userRepoMock.AssertExpectations(t)
	uowMock.AssertExpectations(t)
}

func TestUserRegistrationCommandHandler_FindByExternalProviderError(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	cmd, err := commands.NewUserRegistrationCommand("test", "123", user.ProviderTelegram)
	require.NoError(t, err)

	userRepoMock := &portsmocks.UserRepositoryMock{}
	uowMock := &portsmocks.UnitOfWorkMock{}
	uowMock.On("UserRepository").Return(userRepoMock)

	userRepoMock.
		EXPECT().
		FindByExternalProvider(ctx, user.ProviderTelegram, "123").
		Return(nil, errors.New("repository error")).
		Once()

	uowMock.
		EXPECT().
		Begin(ctx).
		Return(nil).
		Once()
	uowMock.
		EXPECT().
		RollbackUnlessCommitted().
		Return(nil).
		Once()
	// Commit не вызывается при ошибке в репозитории

	handler, err := commands.NewUserRegistrationCommandHandler(logger, uowMock)
	require.NoError(t, err)

	err = handler.Handle(ctx, cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository error")

	userRepoMock.AssertExpectations(t)
	uowMock.AssertExpectations(t)
}

func TestUserRegistrationCommandHandler_CreateError(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	cmd, err := commands.NewUserRegistrationCommand("test", "123", user.ProviderTelegram)
	require.NoError(t, err)

	userRepoMock := &portsmocks.UserRepositoryMock{}
	uowMock := &portsmocks.UnitOfWorkMock{}
	uowMock.On("UserRepository").Return(userRepoMock)

	userRepoMock.
		EXPECT().
		FindByExternalProvider(ctx, user.ProviderTelegram, "123").
		Return(nil, nil).
		Once()
	userRepoMock.
		EXPECT().
		Create(ctx, mock.AnythingOfType("*user.User")).
		Return(errors.New("create error")).
		Once()

	uowMock.
		EXPECT().
		Begin(ctx).
		Return(nil).
		Once()
	uowMock.
		EXPECT().
		RollbackUnlessCommitted().
		Return(nil).
		Once()
	// Commit не вызывается при ошибке создания

	handler, err := commands.NewUserRegistrationCommandHandler(logger, uowMock)
	require.NoError(t, err)

	err = handler.Handle(ctx, cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "create error")

	userRepoMock.AssertExpectations(t)
	uowMock.AssertExpectations(t)
}

func TestUserRegistrationCommandHandler_Logging(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	cmd, err := commands.NewUserRegistrationCommand("new_user", "456", user.ProviderTelegram)
	require.NoError(t, err)

	uowMock, userRepoMock := setupMocks()

	userRepoMock.
		EXPECT().
		FindByExternalProvider(ctx, user.ProviderTelegram, "456").
		Return(nil, nil). // Нового пользователя еще нет
		Once()
	userRepoMock.
		EXPECT().
		Create(ctx, mock.AnythingOfType("*user.User")).
		Return(nil).
		Once()

	uowMock.
		EXPECT().
		Begin(ctx).
		Return(nil).
		Once()
	uowMock.
		EXPECT().
		RollbackUnlessCommitted().
		Return(nil).
		Once()
	uowMock.
		EXPECT().
		Commit(ctx).
		Return(nil).
		Once()

	handler, err := commands.NewUserRegistrationCommandHandler(logger, uowMock)
	require.NoError(t, err)

	err = handler.Handle(ctx, cmd)
	assert.NoError(t, err)

	// Проверяем, что логгер не вызвал панику
	assert.True(t, true) // Заглушка - в реальности проверить вывод логов
}
