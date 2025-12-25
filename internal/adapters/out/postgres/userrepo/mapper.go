package userrepo

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
}

// type ExternalProviderModel
