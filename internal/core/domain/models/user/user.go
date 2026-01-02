package user

import (
	"strings"
	"time"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/ddd"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

type User struct {
	baseAggregate    *ddd.BaseAggregate[shared.ID]
	createdAt        time.Time
	name             string
	externalIdentity *ExternalIdentity
}

func New(name string, chatID string, provider Provider) (*User, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errs.NewValueIsRequiredError("name")
	}

	if chatID == "" || chatID == "0" {
		return nil, errs.NewValueIsRequiredError("chatID")
	}

	u := User{
		baseAggregate: ddd.NewBaseAggregate(shared.NewID()),
		createdAt:     time.Now(),
		name:          name,
	}

	ei, err := NewExternalIdentity(u.ID(), provider, chatID)
	if err != nil {
		return nil, err
	}

	err = u.AddExternalIdentity(ei)
	if err != nil {
		return nil, err
	}

	u.RaiseDomainEvent(NewRegisterEvent(&u))

	return &u, nil
}

func Restore(id shared.ID, name string, createdAt time.Time) *User {
	return &User{
		baseAggregate: ddd.NewBaseAggregate(id),
		name:          name,
		createdAt:     createdAt,
	}
}

func (u *User) AddExternalIdentity(externalIdentity *ExternalIdentity) error {
	if u.externalIdentity != nil {
		return errs.NewValueIsInvalidError("externalIdentity")
	}

	if externalIdentity == nil {
		return errs.NewValueIsRequiredError("externalIdentity")
	}

	if externalIdentity.UserID() != u.ID() {
		return errs.NewValueIsInvalidError("externalIdentity.UserID()")
	}

	u.externalIdentity = externalIdentity

	return nil
}

func (u *User) ID() shared.ID {
	return u.baseAggregate.ID()
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) Name() string {
	return u.name
}

func (u *User) GetExternalIdentity() *ExternalIdentity {
	return u.externalIdentity
}

func (u *User) Equals(other *User) bool {
	if other == nil {
		return false
	}

	return u.baseAggregate.Equal(other.baseAggregate)
}

func (u *User) ClearDomainEvents() {
	u.baseAggregate.ClearDomainEvents()
}

func (u *User) GetDomainEvents() []ddd.DomainEvent {
	return u.baseAggregate.GetDomainEvents()
}

func (u *User) RaiseDomainEvent(event ddd.DomainEvent) {
	u.baseAggregate.RaiseDomainEvent(event)
}
