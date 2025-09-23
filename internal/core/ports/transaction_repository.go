// Package ports определяет интерфейсы портов для взаимодействия с внешними системами.
// Содержит контракты для репозиториев и других внешних зависимостей.
package ports

import (
	"context"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/transaction"
)

// TransactionRepository определяет контракт для работы с хранилищем транзакций.
// Предоставляет методы для добавления, получения, обновления и удаления транзакций.
type TransactionRepository interface {
	// Add добавляет новую транзакцию в хранилище.
	// Возвращает ошибку, если не удалось добавить транзакцию.
	Add(ctx context.Context, transaction *transaction.Transaction) error

	// Get возвращает транзакцию по её идентификатору.
	// Возвращает nil, если транзакция не найдена.
	Get(ctx context.Context, id shared.ID) (*transaction.Transaction, error)

	// Update обновляет существующую транзакцию в хранилище.
	// Возвращает ошибку, если транзакция не найдена или произошла ошибка при обновлении.
	Update(ctx context.Context, transaction *transaction.Transaction) error

	// Delete удаляет транзакцию с указанным идентификатором.
	// Возвращает ошибку, если транзакция не найдена или произошла ошибка при удалении.
	Delete(ctx context.Context, id shared.ID) error
}
