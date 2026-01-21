package commands_test

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Nemizar/coin_tamer_bot/internal/core/application/usecases/commands"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/category"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/user"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
	"github.com/Nemizar/coin_tamer_bot/mocks/core/portsmocks"
)

// setupMocks создает и настраивает моки для тестов
func setupCreateDefaultCategoryMocks() (*portsmocks.UnitOfWorkMock, *portsmocks.UserRepositoryMock, *portsmocks.CategoryRepositoryMock) {
	userRepoMock := &portsmocks.UserRepositoryMock{}
	categoryRepoMock := &portsmocks.CategoryRepositoryMock{}
	uowMock := &portsmocks.UnitOfWorkMock{}
	uowMock.On("UserRepository").Return(userRepoMock)
	uowMock.On("CategoryRepository").Return(categoryRepoMock)
	return uowMock, userRepoMock, categoryRepoMock
}

func TestCreateDefaultCategoryCommandHandler_Validation(t *testing.T) {
	tests := []struct {
		name     string
		chatID   string
		provider user.Provider
		wantErr  error
	}{
		{
			name:     "Empty chat ID",
			chatID:   "",
			provider: user.ProviderTelegram,
			wantErr:  errs.ErrValueIsRequired,
		},
		{
			name:     "Zero chat ID",
			chatID:   "0",
			provider: user.ProviderTelegram,
			wantErr:  errs.ErrValueIsRequired, // Как и в регистрации пользователя, "0" недопустимо
		},
		{
			name:     "Valid data",
			chatID:   "123",
			provider: user.ProviderTelegram,
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := commands.NewCreateDefaultCategoryCommand(tt.chatID, tt.provider)

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

func TestCreateDefaultCategoryCommandHandler_Success(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	// Создаем тестового пользователя
	testUser, err := user.New("test_user", "123", user.ProviderTelegram)
	require.NoError(t, err)

	cmd, err := commands.NewCreateDefaultCategoryCommand("123", user.ProviderTelegram)
	require.NoError(t, err)

	uowMock, userRepoMock, categoryRepoMock := setupCreateDefaultCategoryMocks()

	// Ожидаем, что пользователь будет найден
	userRepoMock.
		EXPECT().
		FindByExternalProvider(ctx, user.ProviderTelegram, "123").
		Return(testUser, nil).
		Once()

	// Ожидаем создание категорий - всего 36 категорий (12 родительских + 24 дочерних)
	categoryRepoMock.
		EXPECT().
		Create(ctx, mock.AnythingOfType("*category.Category")).
		Times(36). // Количество создаваемых категорий: 12 родительских + 24 дочерних
		Return(nil)

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

	handler, err := commands.NewCreateDefaultCategoryCommandHandler(logger, uowMock)
	require.NoError(t, err)

	err = handler.Handle(ctx, cmd)
	assert.NoError(t, err)

	// Проверяем, что моки вызваны в соответствии с ожиданиями
	userRepoMock.AssertExpectations(t)
	categoryRepoMock.AssertExpectations(t)
	uowMock.AssertExpectations(t)
}

func TestCreateDefaultCategoryCommandHandler_UserNotFoundError(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	cmd, err := commands.NewCreateDefaultCategoryCommand("123", user.ProviderTelegram)
	require.NoError(t, err)

	userRepoMock := &portsmocks.UserRepositoryMock{}
	categoryRepoMock := &portsmocks.CategoryRepositoryMock{}
	uowMock := &portsmocks.UnitOfWorkMock{}
	uowMock.On("UserRepository").Return(userRepoMock)
	uowMock.On("CategoryRepository").Return(categoryRepoMock).Maybe() // Указываем, что вызов может не произойти

	// Ожидаем, что пользователь не будет найден
	userRepoMock.
		EXPECT().
		FindByExternalProvider(ctx, user.ProviderTelegram, "123").
		Return(nil, errs.ErrObjectNotFound).
		Once()

	// Create не должен вызываться, так как пользователь не найден
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
	// Commit не вызывается при ошибке поиска пользователя

	handler, err := commands.NewCreateDefaultCategoryCommandHandler(logger, uowMock)
	require.NoError(t, err)

	err = handler.Handle(ctx, cmd)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrObjectNotFound))

	// Проверяем, что Create не был вызван
	categoryRepoMock.AssertNotCalled(t, "Create")
	userRepoMock.AssertExpectations(t)
	uowMock.AssertExpectations(t)
}

func TestCreateDefaultCategoryCommandHandler_BeginTransactionError(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	cmd, err := commands.NewCreateDefaultCategoryCommand("123", user.ProviderTelegram)
	require.NoError(t, err)

	// В этом тесте репозитории не будут вызваны, так как ошибка происходит в Begin
	userRepoMock := &portsmocks.UserRepositoryMock{}
	categoryRepoMock := &portsmocks.CategoryRepositoryMock{}
	uowMock := &portsmocks.UnitOfWorkMock{}
	uowMock.On("UserRepository").Return(userRepoMock).Maybe()
	uowMock.On("CategoryRepository").Return(categoryRepoMock).Maybe()

	// Мок возвращает ошибку при начале транзакции
	uowMock.
		EXPECT().
		Begin(ctx).
		Return(errors.New("begin transaction error")).
		Once()
	uowMock.
		EXPECT().
		RollbackUnlessCommitted().
		Return(nil).
		Once()

	handler, err := commands.NewCreateDefaultCategoryCommandHandler(logger, uowMock)
	require.NoError(t, err)

	err = handler.Handle(ctx, cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "begin transaction error")

	uowMock.AssertExpectations(t)
}

func TestCreateDefaultCategoryCommandHandler_CommitError(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	// Создаем тестового пользователя
	testUser, err := user.New("test_user", "123", user.ProviderTelegram)
	require.NoError(t, err)

	cmd, err := commands.NewCreateDefaultCategoryCommand("123", user.ProviderTelegram)
	require.NoError(t, err)

	uowMock, userRepoMock, categoryRepoMock := setupCreateDefaultCategoryMocks()

	// Ожидаем, что пользователь будет найден
	userRepoMock.
		EXPECT().
		FindByExternalProvider(ctx, user.ProviderTelegram, "123").
		Return(testUser, nil).
		Once()

	// Ожидаем создание категорий
	categoryRepoMock.
		EXPECT().
		Create(ctx, mock.AnythingOfType("*category.Category")).
		Times(36). // Количество создаваемых категорий: 12 родительских + 24 дочерних
		Return(nil)

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

	handler, err := commands.NewCreateDefaultCategoryCommandHandler(logger, uowMock)
	require.NoError(t, err)

	err = handler.Handle(ctx, cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "commit error")

	userRepoMock.AssertExpectations(t)
	categoryRepoMock.AssertExpectations(t)
	uowMock.AssertExpectations(t)
}

func TestCreateDefaultCategoryCommandHandler_FindUserError(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	cmd, err := commands.NewCreateDefaultCategoryCommand("123", user.ProviderTelegram)
	require.NoError(t, err)

	userRepoMock := &portsmocks.UserRepositoryMock{}
	categoryRepoMock := &portsmocks.CategoryRepositoryMock{}
	uowMock := &portsmocks.UnitOfWorkMock{}
	uowMock.On("UserRepository").Return(userRepoMock)
	uowMock.On("CategoryRepository").Return(categoryRepoMock).Maybe() // Указываем, что вызов может не произойти

	// Ожидаем, что поиск пользователя вернет ошибку
	findError := errors.New("find user error")
	userRepoMock.
		EXPECT().
		FindByExternalProvider(ctx, user.ProviderTelegram, "123").
		Return(nil, findError).
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
	// Commit не вызывается при ошибке поиска пользователя

	handler, err := commands.NewCreateDefaultCategoryCommandHandler(logger, uowMock)
	require.NoError(t, err)

	err = handler.Handle(ctx, cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "find user error")

	userRepoMock.AssertExpectations(t)
	uowMock.AssertExpectations(t)
	// Create не должен быть вызван
	categoryRepoMock.AssertNotCalled(t, "Create")
}

func TestCreateDefaultCategoryCommandHandler_CategoryCreateError(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	// Создаем тестового пользователя
	testUser, err := user.New("test_user", "123", user.ProviderTelegram)
	require.NoError(t, err)

	cmd, err := commands.NewCreateDefaultCategoryCommand("123", user.ProviderTelegram)
	require.NoError(t, err)

	uowMock, userRepoMock, categoryRepoMock := setupCreateDefaultCategoryMocks()

	// Ожидаем, что пользователь будет найден
	userRepoMock.
		EXPECT().
		FindByExternalProvider(ctx, user.ProviderTelegram, "123").
		Return(testUser, nil).
		Once()

	// Ожидаем, что создание категории вернет ошибку
	createError := errors.New("create category error")
	categoryRepoMock.
		EXPECT().
		Create(ctx, mock.AnythingOfType("*category.Category")).
		Return(createError).
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
	// Commit не вызывается при ошибке создания категории

	handler, err := commands.NewCreateDefaultCategoryCommandHandler(logger, uowMock)
	require.NoError(t, err)

	err = handler.Handle(ctx, cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "create category error")

	userRepoMock.AssertExpectations(t)
	categoryRepoMock.AssertExpectations(t)
	uowMock.AssertExpectations(t)
}

func TestCreateDefaultCategoryCommandHandler_CategoryTypeValidation(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	// Создаем тестового пользователя
	testUser, err := user.New("test_user", "123", user.ProviderTelegram)
	require.NoError(t, err)

	cmd, err := commands.NewCreateDefaultCategoryCommand("123", user.ProviderTelegram)
	require.NoError(t, err)

	uowMock, userRepoMock, categoryRepoMock := setupCreateDefaultCategoryMocks()

	// Ожидаем, что пользователь будет найден
	userRepoMock.
		EXPECT().
		FindByExternalProvider(ctx, user.ProviderTelegram, "123").
		Return(testUser, nil).
		Once()

	// Проверяем, что создаваемые категории имеют правильный тип
	categoryRepoMock.
		EXPECT().
		Create(ctx, mock.MatchedBy(func(cat *category.Category) bool {
			// Проверяем, что категория имеет правильный тип (income или expense)
			return cat.Type() == category.TypeIncome || cat.Type() == category.TypeExpense
		})).
		Times(36). // Количество создаваемых категорий: 12 родительских + 24 дочерних
		Return(nil)

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

	handler, err := commands.NewCreateDefaultCategoryCommandHandler(logger, uowMock)
	require.NoError(t, err)

	err = handler.Handle(ctx, cmd)
	assert.NoError(t, err)

	// Проверяем, что моки вызваны в соответствии с ожиданиями
	userRepoMock.AssertExpectations(t)
	categoryRepoMock.AssertExpectations(t)
	uowMock.AssertExpectations(t)
}
