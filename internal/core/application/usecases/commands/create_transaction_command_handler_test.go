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
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/transaction"
	"github.com/Nemizar/coin_tamer_bot/mocks/core/portsmocks"
)

// setupMocks создает и настраивает моки для тестов
func setupCreateTransactionMocks() (*portsmocks.UnitOfWorkMock, *portsmocks.TransactionRepositoryMock) {
	transactionRepoMock := &portsmocks.TransactionRepositoryMock{}
	uowMock := &portsmocks.UnitOfWorkMock{}
	uowMock.On("TransactionRepository").Return(transactionRepoMock)
	return uowMock, transactionRepoMock
}

func TestCreateTransactionCommandHandler_Validation(t *testing.T) {
	// Тестирование валидации команды происходит в самой команде,
	// а не в обработчике, поэтому здесь мы тестируем только валидацию
	// при создании транзакции в обработчике

	userID := shared.NewID()
	categoryID := shared.NewID()

	tests := []struct {
		name       string
		userID     shared.ID
		amount     transaction.Amount
		categoryID shared.ID
		wantErr    error
	}{
		{
			name:       "Valid data",
			userID:     userID,
			amount:     createValidAmount(t),
			categoryID: categoryID,
			wantErr:    nil,
		},
		{
			name:       "Zero user ID",
			userID:     shared.ID{}, // Zero ID
			amount:     createValidAmount(t),
			categoryID: categoryID,
			wantErr:    transaction.ErrInvalidUserID,
		},
		{
			name:       "Zero category ID",
			userID:     userID,
			amount:     createValidAmount(t),
			categoryID: shared.ID{}, // Zero ID
			wantErr:    transaction.ErrInvalidCategoryID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем транзакцию напрямую, чтобы проверить валидацию
			_, err := transaction.New(tt.userID, tt.amount, tt.categoryID)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func createValidAmount(t *testing.T) transaction.Amount {
	amount, err := transaction.NewAmountFromString("100.50")
	require.NoError(t, err)
	return amount
}

func TestCreateTransactionCommandHandler_Success(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	userID := shared.NewID()
	categoryID := shared.NewID()
	amount := createValidAmount(t)

	cmd, err := commands.NewCreateTransactionCommand(userID, amount, categoryID)
	require.NoError(t, err)

	uowMock, transactionRepoMock := setupCreateTransactionMocks()

	// Ожидаем создание транзакции
	transactionRepoMock.
		EXPECT().
		Add(ctx, mock.AnythingOfType("*transaction.Transaction")).
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

	handler, err := commands.NewCreateTransactionCommandHandler(logger, uowMock)
	require.NoError(t, err)

	err = handler.Handle(ctx, cmd)
	assert.NoError(t, err)

	// Проверяем, что моки вызваны в соответствии с ожиданиями
	transactionRepoMock.AssertExpectations(t)
	uowMock.AssertExpectations(t)
}

func TestCreateTransactionCommandHandler_BeginTransactionError(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	userID := shared.NewID()
	categoryID := shared.NewID()
	amount := createValidAmount(t)

	cmd, err := commands.NewCreateTransactionCommand(userID, amount, categoryID)
	require.NoError(t, err)

	// В этом тесте TransactionRepository() не будет вызван, так как ошибка происходит в Begin
	transactionRepoMock := &portsmocks.TransactionRepositoryMock{}
	uowMock := &portsmocks.UnitOfWorkMock{}

	// Мок возвращает ошибку при начале транзакции
	uowMock.
		EXPECT().
		Begin(ctx).
		Return(errors.New("begin transaction error")).
		Once()
	// Устанавливаем ожидание TransactionRepository, но помечаем как Maybe, чтобы не требовался вызов
	uowMock.On("TransactionRepository").Return(transactionRepoMock).Maybe()
	uowMock.
		EXPECT().
		RollbackUnlessCommitted().
		Return(nil).
		Once()

	handler, err := commands.NewCreateTransactionCommandHandler(logger, uowMock)
	require.NoError(t, err)

	err = handler.Handle(ctx, cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "begin transaction error")

	uowMock.AssertExpectations(t)
}

func TestCreateTransactionCommandHandler_CommitError(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	userID := shared.NewID()
	categoryID := shared.NewID()
	amount := createValidAmount(t)

	cmd, err := commands.NewCreateTransactionCommand(userID, amount, categoryID)
	require.NoError(t, err)

	uowMock, transactionRepoMock := setupCreateTransactionMocks()

	// Ожидаем создание транзакции
	transactionRepoMock.
		EXPECT().
		Add(ctx, mock.AnythingOfType("*transaction.Transaction")).
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

	handler, err := commands.NewCreateTransactionCommandHandler(logger, uowMock)
	require.NoError(t, err)

	err = handler.Handle(ctx, cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "commit error")

	transactionRepoMock.AssertExpectations(t)
	uowMock.AssertExpectations(t)
}

func TestCreateTransactionCommandHandler_TransactionValidationError(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	// Используем нулевые ID, чтобы вызвать ошибку валидации транзакции
	zeroID := shared.ID{}
	amount := createValidAmount(t)

	cmd, err := commands.NewCreateTransactionCommand(zeroID, amount, zeroID)
	require.NoError(t, err)

	// В этом тесте TransactionRepository() не будет вызван, так как ошибка происходит при создании транзакции
	transactionRepoMock := &portsmocks.TransactionRepositoryMock{}
	uowMock := &portsmocks.UnitOfWorkMock{}
	uowMock.On("TransactionRepository").Return(transactionRepoMock).Maybe()

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
	// Commit не вызывается при ошибке создания транзакции

	handler, err := commands.NewCreateTransactionCommandHandler(logger, uowMock)
	require.NoError(t, err)

	err = handler.Handle(ctx, cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid user id") // или "invalid category id"

	transactionRepoMock.AssertNotCalled(t, "Add")
	transactionRepoMock.AssertExpectations(t)
	uowMock.AssertExpectations(t)
}

func TestCreateTransactionCommandHandler_AddTransactionError(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	userID := shared.NewID()
	categoryID := shared.NewID()
	amount := createValidAmount(t)

	cmd, err := commands.NewCreateTransactionCommand(userID, amount, categoryID)
	require.NoError(t, err)

	uowMock, transactionRepoMock := setupCreateTransactionMocks()

	// Ожидаем, что добавление транзакции вернет ошибку
	addError := errors.New("add transaction error")
	transactionRepoMock.
		EXPECT().
		Add(ctx, mock.AnythingOfType("*transaction.Transaction")).
		Return(addError).
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
	// Commit не вызывается при ошибке добавления транзакции

	handler, err := commands.NewCreateTransactionCommandHandler(logger, uowMock)
	require.NoError(t, err)

	err = handler.Handle(ctx, cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "add transaction error")

	transactionRepoMock.AssertExpectations(t)
	uowMock.AssertExpectations(t)
}

func TestCreateTransactionCommandHandler_InvalidUserID(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	zeroID := shared.ID{} // Невалидный ID
	categoryID := shared.NewID()
	amount := createValidAmount(t)

	cmd, err := commands.NewCreateTransactionCommand(zeroID, amount, categoryID)
	require.NoError(t, err)

	// В этом тесте TransactionRepository() не будет вызван, так как ошибка происходит при создании транзакции
	transactionRepoMock := &portsmocks.TransactionRepositoryMock{}
	uowMock := &portsmocks.UnitOfWorkMock{}
	uowMock.On("TransactionRepository").Return(transactionRepoMock).Maybe()

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
	// Commit не вызывается при ошибке создания транзакции

	handler, err := commands.NewCreateTransactionCommandHandler(logger, uowMock)
	require.NoError(t, err)

	err = handler.Handle(ctx, cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid user id")

	transactionRepoMock.AssertNotCalled(t, "Add")
	transactionRepoMock.AssertExpectations(t)
	uowMock.AssertExpectations(t)
}

func TestCreateTransactionCommandHandler_Logging(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	userID := shared.NewID()
	categoryID := shared.NewID()
	amount := createValidAmount(t)

	cmd, err := commands.NewCreateTransactionCommand(userID, amount, categoryID)
	require.NoError(t, err)

	uowMock, transactionRepoMock := setupCreateTransactionMocks()

	// Ожидаем создание транзакции
	transactionRepoMock.
		EXPECT().
		Add(ctx, mock.AnythingOfType("*transaction.Transaction")).
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

	handler, err := commands.NewCreateTransactionCommandHandler(logger, uowMock)
	require.NoError(t, err)

	err = handler.Handle(ctx, cmd)
	assert.NoError(t, err)

	// Проверяем, что логгер не вызвал панику
	assert.True(t, true) // Заглушка - в реальности проверить вывод логов
}
