// Package ports определяет интерфейсы портов для взаимодействия с внешними системами.
// Содержит контракты для работы с единицами работы (Unit of Work).
package ports

// UnitOfWork определяет интерфейс для управления транзакциями и доступа к репозиториям.
// Обеспечивает атомарность операций и согласованность данных.
type UnitOfWork interface {
	// Begin начинает новую транзакцию.
	// Все последующие операции будут выполнены в рамках этой транзакции.
	Begin() error

	// Commit фиксирует текущую транзакцию.
	// Возвращает ошибку, если не удалось зафиксировать изменения.
	Commit() error

	UserRepository() UserRepository
	CategoryRepository() CategoryRepository
	TransactionRepository() TransactionRepository
	ExternalIdentityRepository() ExternalIdentityRepository

	// RollbackUnlessCommitted откатывает текущую транзакцию, если она не была зафиксирована.
	// Используется для безопасной отмены изменений в случае ошибки.
	RollbackUnlessCommitted() error
}
