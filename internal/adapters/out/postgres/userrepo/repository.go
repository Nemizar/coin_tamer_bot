package userrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/identity"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/user"
	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

type UserRepository struct {
	tracker Tracker
}

func NewUserRepository(tracker Tracker) ports.UserRepository {
	return &UserRepository{tracker: tracker}
}

func (u UserRepository) Create(ctx context.Context, user *user.User) error {
	u.tracker.Track(user)

	stmt := `INSERT INTO users (id, name, created_at)
			 VALUES ($1, $2, $3)`
	_, err := u.tracker.DB().ExecContext(ctx, stmt, user.ID(), user.Name(), user.CreatedAt())
	if err != nil {
		return fmt.Errorf("user repo insert: %w", err)
	}

	return nil
}

func (u UserRepository) FindByExternalProvider(provider identity.Provider, externalID string) (*user.User, error) {
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
