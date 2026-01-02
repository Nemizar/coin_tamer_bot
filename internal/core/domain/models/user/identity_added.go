package user

import (
	"reflect"

	"github.com/google/uuid"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"

	"github.com/Nemizar/coin_tamer_bot/internal/pkg/ddd"
)

var _ ddd.DomainEvent = &RegisterEvent{}

type RegisterEvent struct {
	id   uuid.UUID
	name string

	payload dto

	isValid bool
}

func NewRegisterEvent(u *User) *RegisterEvent {
	d := dto{
		ID:         u.GetExternalIdentity().ID().Value(),
		UserID:     u.GetExternalIdentity().UserID(),
		UserName:   u.Name(),
		Provider:   u.GetExternalIdentity().Provider(),
		ExternalID: u.GetExternalIdentity().ExternalID(),
	}

	event := RegisterEvent{
		id:      uuid.New(),
		payload: d,
		isValid: true,
	}

	event.name = reflect.TypeOf(event).Name()

	return &event
}

func NewEmptyRegisterEvent() *RegisterEvent {
	event := RegisterEvent{}
	event.name = reflect.TypeOf(event).Name()
	return &event
}

func (c RegisterEvent) GetID() uuid.UUID {
	return c.id
}

func (c RegisterEvent) GetName() string {
	return c.name
}

func (c RegisterEvent) GetPayload() dto {
	return c.payload
}

func (c RegisterEvent) IsValid() bool {
	return c.isValid
}

type dto struct {
	ID         uuid.UUID
	UserID     shared.ID
	UserName   string
	Provider   Provider
	ExternalID string
}
