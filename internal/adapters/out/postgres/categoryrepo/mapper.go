package categoryrepo

import (
	"time"

	"github.com/google/uuid"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/category"
)

type Model struct {
	ID           uuid.UUID
	Name         string
	OwnerID      uuid.UUID
	ParentID     uuid.UUID
	CategoryType category.Type
	CreatedAt    time.Time
}
