package transactionrepo

import (
	"github.com/jmoiron/sqlx"

	"github.com/Nemizar/coin_tamer_bot/internal/pkg/ddd"
)

type Tracker interface {
	Tx() *sqlx.Tx
	DB() *sqlx.DB
	InTx() bool
	Track(agg ddd.AggregateRoot)
	Begin() error
	Commit() error
}
