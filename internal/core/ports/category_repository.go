// Package ports определяет интерфейсы портов для взаимодействия с внешними системами.
// Содержит контракты для репозиториев и других внешних зависимостей.
package ports

import (
	"context"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/category"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
)

// CategoryRepository определяет контракт для работы с хранилищем категорий.
// Предоставляет методы для создания, чтения, обновления и удаления категорий.
type CategoryRepository interface {
	// Create сохраняет новую категорию в хранилище.
	// Возвращает ошибку, если не удалось создать категорию.
	Create(ctx context.Context, category *category.Category) error

	// GetAll возвращает все категории из хранилища.
	// Возвращает пустой список, если категории не найдены.
	GetAll(ctx context.Context) ([]*category.Category, error)

	// Update обновляет существующую категорию в хранилище.
	// Возвращает ошибку, если категория не найдена или произошла ошибка при обновлении.
	Update(ctx context.Context, category *category.Category) error

	// Delete удаляет категорию с указанным идентификатором.
	// Возвращает ошибку, если категория не найдена или произошла ошибка при удалении.
	Delete(ctx context.Context, id shared.ID) error
}
