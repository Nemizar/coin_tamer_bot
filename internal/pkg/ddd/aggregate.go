// Package ddd содержит базовые структуры и интерфейсы для реализации предметно-ориентированного дизайна (DDD).
// Определяет базовые строительные блоки для доменных моделей, такие как агрегаты и сущности.
package ddd

// BaseAggregate представляет базовый агрегат, который инкапсулирует бизнес-правила и обеспечивает целостность данных.
// Содержит идентификатор и список доменных событий.
type BaseAggregate[ID comparable] struct {
	baseEntity   *BaseEntity[ID]
	domainEvents []DomainEvent
}

// NewBaseAggregate создает новый экземпляр BaseAggregate с указанным идентификатором.
// Инициализирует пустой список доменных событий.
func NewBaseAggregate[ID comparable](id ID) *BaseAggregate[ID] {
	return &BaseAggregate[ID]{
		baseEntity:   NewBaseEntity[ID](id),
		domainEvents: make([]DomainEvent, 0),
	}
}

// ID возвращает уникальный идентификатор агрегата.
func (ba *BaseAggregate[ID]) ID() ID {
	return ba.baseEntity.ID()
}

// Equal сравнивает текущий агрегат с другим агрегатом на основе их идентификаторов.
// Возвращает true, если идентификаторы совпадают, иначе false.
func (ba *BaseAggregate[ID]) Equal(other *BaseAggregate[ID]) bool {
	if other == nil {
		return false
	}
	return ba.baseEntity.Equal(other.baseEntity)
}

// ClearDomainEvents очищает список доменных событий агрегата.
// Используется после успешной обработки событий.
func (ba *BaseAggregate[ID]) ClearDomainEvents() {
	ba.domainEvents = []DomainEvent{}
}

// GetDomainEvents возвращает список всех доменных событий, произошедших с агрегатом.
// Возвращает копию списка событий для предотвращения нежелательных изменений.
func (ba *BaseAggregate[ID]) GetDomainEvents() []DomainEvent {
	return ba.domainEvents
}

// RaiseDomainEvent добавляет новое доменное событие в список событий агрегата.
// Используется для регистрации изменений состояния агрегата.
func (ba *BaseAggregate[ID]) RaiseDomainEvent(event DomainEvent) {
	ba.domainEvents = append(ba.domainEvents, event)
}
