package ports

type UnitOfWorkFactory interface {
	New() (UnitOfWork, error)
}
