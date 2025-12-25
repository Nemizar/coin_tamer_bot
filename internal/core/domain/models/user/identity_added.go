package user

import (
	"reflect"

	"github.com/google/uuid"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"

	"github.com/Nemizar/coin_tamer_bot/internal/pkg/ddd"
)

var _ ddd.DomainEvent = &ExternalIdentityAddedEvent{}

type ExternalIdentityAddedEvent struct {
	id   uuid.UUID
	name string

	payload dto

	isValid bool
}

func NewExternalIdentityAddedEvent(identity *ExternalIdentity) *ExternalIdentityAddedEvent {
	d := dto{
		ID:         identity.ID().Value(),
		UserID:     identity.UserID(),
		Provider:   identity.Provider(),
		ExternalID: identity.ExternalID(),
	}

	event := ExternalIdentityAddedEvent{
		id:      uuid.New(),
		payload: d,
		isValid: true,
	}

	event.name = reflect.TypeOf(event).Name()

	return &event
}

func NewEmptyExternalIdentityAddedEvent() *ExternalIdentityAddedEvent {
	event := ExternalIdentityAddedEvent{}
	event.name = reflect.TypeOf(event).Name()
	return &event
}

func (c ExternalIdentityAddedEvent) GetID() uuid.UUID {
	return c.id
}

func (c ExternalIdentityAddedEvent) GetName() string {
	return c.name
}

func (c ExternalIdentityAddedEvent) GetPayload() interface{} {
	return c.payload
}

func (c ExternalIdentityAddedEvent) IsValid() bool {
	return c.isValid
}

type dto struct {
	ID         uuid.UUID
	UserID     shared.ID
	Provider   Provider
	ExternalID string
}
