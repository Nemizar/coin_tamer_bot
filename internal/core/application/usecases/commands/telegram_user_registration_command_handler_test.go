package commands_test

import (
	"context"
	"log"
	"testing"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/identity"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"

	"github.com/Nemizar/coin_tamer_bot/internal/adapters/out/postgres"
	"github.com/Nemizar/coin_tamer_bot/internal/core/application/usecases/commands"
	"github.com/Nemizar/coin_tamer_bot/internal/migrations"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/testcnts"
)

func setupTest(t *testing.T) (context.Context, *sqlx.DB) {
	ctx := context.Background()

	postgresContainer, dsn, err := testcnts.StartPostgresContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	pool, err := sqlx.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}

	goose.SetBaseFS(migrations.FS)
	if err := goose.Up(pool.DB, "."); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	t.Cleanup(func() {
		err := pool.Close()
		assert.NoError(t, err)

		err = postgresContainer.Terminate(ctx)
		assert.NoError(t, err)
	})
	return ctx, pool
}

func TestUserRegistrationCommandHandler_Success(t *testing.T) {
	ctx, pool := setupTest(t)

	uowFactory, err := postgres.NewUnitOfWorkFactory(pool)
	require.Nil(t, err)

	cmd, err := commands.NewUserRegistrationCommand("test", 123)
	assert.Nil(t, err)

	handler := commands.NewUserRegistrationCommandHandler(nil, uowFactory)

	err = handler.Handle(ctx, cmd)
	assert.Nil(t, err)

	ei, err := uowFactory.New(ctx)
	require.Nil(t, err)

	u, err := ei.UserRepository().FindByExternalProvider(identity.ProviderTelegram, "123")
	require.Nil(t, err)
	require.NotNil(t, u)
	assert.Equal(t, "test", u.Name())
	assert.NotEqual(t, uuid.Nil, u.ID())
}

func TestUserRegistrationCommandHandler_Failure_EmptyName(t *testing.T) {
	cmd, err := commands.NewUserRegistrationCommand("", 123)
	assert.Nil(t, cmd)
	assert.ErrorIs(t, err, errs.ErrValueIsRequired)
}

func TestUserRegistrationCommandHandler_Failure_EmptyTelegramChatID(t *testing.T) {
	cmd, err := commands.NewUserRegistrationCommand("test", 0)
	assert.Nil(t, cmd)
	assert.ErrorIs(t, err, errs.ErrValueIsRequired)
}
