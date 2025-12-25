package userrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/user"
	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

type UserRepository struct {
	tracker Tracker
}

func NewUserRepository(tracker Tracker) (ports.UserRepository, error) {
	if tracker == nil {
		return nil, errs.NewValueIsRequiredError("tracker")
	}

	return &UserRepository{tracker: tracker}, nil
}

func (u UserRepository) Create(ctx context.Context, us *user.User) error {
	u.tracker.Track(us)

	if !u.tracker.InTx() {
		err := u.tracker.Begin(ctx)
		if err != nil {
			return fmt.Errorf("user repo begin: %w", err)
		}

		defer func(tracker Tracker, ctx context.Context) {
			err = tracker.Commit(ctx)
			if err != nil {
				tracker.Logger().Error("user repo commit", "err", err)
			}
		}(u.tracker, ctx)
	}

	stmt := `INSERT INTO users (id, name, created_at)
			 VALUES ($1, $2, $3)`
	_, err := u.tracker.Tx().ExecContext(ctx, stmt, us.ID(), us.Name(), us.CreatedAt())
	if err != nil {
		return fmt.Errorf("user repo insert: %w", err)
	}

	for _, ei := range us.GetExternalIdentities() {
		stmt = `INSERT INTO external_identities (id, user_id, provider, external_id, created_at)
			VALUES ($1, $2, $3, $4, $5)`
		_, err = u.tracker.Tx().ExecContext(ctx, stmt, ei.ID(), ei.UserID(), ei.Provider(), ei.ExternalID(), ei.GetCreatedAt())
		if err != nil {
			return fmt.Errorf("user repo insert external identity: %w", err)
		}
	}

	return nil
}

func (u UserRepository) FindByExternalProvider(provider user.Provider, externalID string) (*user.User, error) {
	stmt := `SELECT u.id, name, u.created_at
				FROM users u
				INNER JOIN external_identities ei ON u.id = ei.user_id
				WHERE ei.external_id = $1 AND ei.provider = $2`
	row := u.tracker.DB().QueryRowContext(context.Background(), stmt, externalID, provider)

	var repoModel Model
	err := row.Scan(&repoModel.ID, &repoModel.Name, &repoModel.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewObjectNotFoundError("user", externalID)
		}

		return nil, fmt.Errorf("user repo find by external provider: %w", err)
	}

	return user.Restore(shared.RestoreID(repoModel.ID), repoModel.Name, repoModel.CreatedAt), nil
}
