// Package ddd содержит базовые структуры и интерфейсы для реализации предметно-ориентированного дизайна (DDD).
// Реализует паттерн Mediator для обработки доменных событий.
package ddd

import "context"

// EventHandler определяет контракт для обработчиков доменных событий.
// Каждый обработчик должен реализовывать метод Handle для обработки конкретного типа события.
// Обработчик должен быть потокобезопасным, так как может вызываться из нескольких горутин.
type EventHandler interface {
	Handle(ctx context.Context, event DomainEvent) error
}

// Mediatr определяет интерфейс для посредника, управляющего подпиской и публикацией событий.
// Позволяет подписывать обработчики на события и публиковать события для обработки.
type Mediatr interface {
	Subscribe(handler EventHandler, events ...DomainEvent)
	Publish(ctx context.Context, event DomainEvent) error
}

// mediatr реализует интерфейс Mediatr.
// Хранит отображение имен событий на список обработчиков.
type mediatr struct {
	handlers map[string][]EventHandler
}

// NewMediatr создает и возвращает новый экземпляр посредника.
// Инициализирует внутреннюю карту для хранения обработчиков событий.
func NewMediatr() Mediatr {
	return &mediatr{handlers: make(map[string][]EventHandler)}
}

// Subscribe регистрирует обработчик для указанных типов событий.
// Один обработчик может быть подписан на несколько типов событий.
// При наступлении события вызываются все зарегистрированные для него обработчики.
func (e *mediatr) Subscribe(handler EventHandler, events ...DomainEvent) {
	for _, event := range events {
		handlers := e.handlers[event.GetName()]
		handlers = append(handlers, handler)
		e.handlers[event.GetName()] = handlers
	}
}

// Publish публикует событие для обработки.
// Вызывает все обработчики, подписанные на данное событие.
// Возвращает первую возникшую ошибку или nil, если обработка прошла успешно.
func (e *mediatr) Publish(ctx context.Context, event DomainEvent) error {
	for _, handler := range e.handlers[event.GetName()] {
		err := handler.Handle(ctx, event)
		if err != nil {
			return err
		}
	}
	return nil
}
