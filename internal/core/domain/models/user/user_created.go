package user

import (
	"reflect"

	"github.com/google/uuid"

	"github.com/Nemizar/coin_tamer_bot/internal/pkg/ddd"
)

var _ ddd.DomainEvent = &CreatedEvent{}

type CreatedEvent struct {
	id   uuid.UUID
	name string

	payload dto

	isValid bool
}

func NewCreateEvent(aggregate *User) *CreatedEvent {
	d := dto{
		ID: aggregate.ID().Value(),
	}

	event := CreatedEvent{
		id:      uuid.New(),
		payload: d,
		isValid: true,
	}

	event.name = reflect.TypeOf(event).Name()

	return &event
}

func NewEmptyCreateEvent() *CreatedEvent {
	event := CreatedEvent{}
	event.name = reflect.TypeOf(event).Name()
	return &event
}

func (c CreatedEvent) GetID() uuid.UUID {
	return c.id
}

func (c CreatedEvent) GetName() string {
	return c.name
}

func (c CreatedEvent) IsValid() bool {
	return c.isValid
}

type dto struct {
	ID uuid.UUID
}
