// Package shared предоставляет общие типы домена и утилиты, используемые в различных пакетах.
// Включает объекты-значения и базовые типы для обеспечения согласованности и уменьшения дублирования кода.
package shared

import "github.com/google/uuid"

// ID представляет объект-значение уникального идентификатора с использованием UUID v7 с fallback на v4.
// Обеспечивает типобезопасную работу с идентификаторами в доменной модели.
type ID struct {
	value uuid.UUID
}

// NewID создает новый экземпляр ID с использованием формата UUID v7 с временным упорядочиванием.
// При ошибке генерации v7 переходит на UUID v4.
func NewID() ID {
	u, err := uuid.NewV7()
	if err != nil {
		u = uuid.New()
	}
	return ID{value: u}
}

// NewIDFromString создает ID из строкового представления UUID.
// Возвращает ошибку если строка не является валидным UUID.
func NewIDFromString(s string) (ID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return ID{}, err
	}
	return ID{value: id}, nil
}

// String возвращает строковое представление ID.
func (id ID) String() string {
	return id.value.String()
}

// IsZero возвращает true если ID не инициализирован (равен uuid.Nil).
func (id ID) IsZero() bool {
	return id.value == uuid.Nil
}

// Value возвращает underlying значение uuid.UUID.
func (id ID) Value() uuid.UUID {
	return id.value
}
