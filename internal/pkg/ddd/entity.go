// Package ddd содержит базовые структуры и интерфейсы для реализации предметно-ориентированного дизайна (DDD).
// Определяет основные строительные блоки для доменных моделей.
package ddd

// BaseEntity представляет базовую сущность с общими полями и методами для всех доменных сущностей.
// Содержит идентификатор, который может быть любого сопоставимого типа.
type BaseEntity[ID comparable] struct {
	id ID
}

// NewBaseEntity создает новый экземпляр BaseEntity с указанным идентификатором.
// Используется как базовая реализация для всех доменных сущностей.
func NewBaseEntity[ID comparable](id ID) *BaseEntity[ID] {
	return &BaseEntity[ID]{
		id: id,
	}
}

// ID возвращает уникальный идентификатор сущности.
func (be *BaseEntity[ID]) ID() ID {
	return be.id
}

// Equal сравнивает текущую сущность с другой сущностью.
// Возвращает true, если сущности имеют одинаковый идентификатор.
func (be *BaseEntity[ID]) Equal(other *BaseEntity[ID]) bool {
	if other == nil {
		return false
	}
	return be.id == other.id
}
