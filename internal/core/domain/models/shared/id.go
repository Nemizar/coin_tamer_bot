package shared

import "github.com/google/uuid"

type ID struct {
	value uuid.UUID
}

func NewID() ID {
	u, err := uuid.NewV7()
	if err != nil {
		u = uuid.New()
	}
	return ID{value: u}
}

func RestoreID(id uuid.UUID) ID {
	return ID{value: id}
}

func NewIDFromString(s string) (ID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return ID{}, err
	}
	return ID{value: id}, nil
}

func (id ID) String() string {
	return id.value.String()
}

func (id ID) IsZero() bool {
	return id.value == uuid.Nil
}

func (id ID) Value() uuid.UUID {
	return id.value
}
