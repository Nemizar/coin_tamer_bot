package categoryrepo

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"

	"github.com/Nemizar/coin_tamer_bot/internal/pkg/ddd"
)

type Tracker interface {
	Tx() *sqlx.Tx
	DB() *sqlx.DB
	InTx() bool
	Track(agg ddd.AggregateRoot)
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Logger() ports.Logger
}
