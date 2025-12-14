package ports

import "context"

type UnitOfWork interface {
	Begin(ctx context.Context) error

	Commit(ctx context.Context) error

	UserRepository() UserRepository
	CategoryRepository() CategoryRepository
	TransactionRepository() TransactionRepository
	ExternalIdentityRepository() ExternalIdentityRepository

	RollbackUnlessCommitted() error
}
