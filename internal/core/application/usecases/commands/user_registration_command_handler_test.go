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

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/user"

	cmd2 "github.com/Nemizar/coin_tamer_bot/cmd"
	"github.com/Nemizar/coin_tamer_bot/configs"

	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"

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

	cr := cmd2.NewCompositionRoot(configs.Config{}, pool)

	uowFactory := cr.NewUnitOfWorkFactory()

	cmd, err := commands.NewUserRegistrationCommand("test", "123", user.ProviderTelegram)
	assert.Nil(t, err)

	handler := cr.NewUserRegistrationCommandHandler()

	err = handler.Handle(ctx, cmd)
	assert.Nil(t, err)

	ei, err := uowFactory.New()
	require.Nil(t, err)

	u, err := ei.UserRepository().FindByExternalProvider(user.ProviderTelegram, "123")
	require.Nil(t, err)
	require.NotNil(t, u)
	assert.Equal(t, "test", u.Name())
	assert.NotEqual(t, uuid.Nil, u.ID())
}

func TestUserRegistrationCommandHandler_Failure_EmptyName(t *testing.T) {
	cmd, err := commands.NewUserRegistrationCommand("", "123", user.ProviderTelegram)
	assert.Nil(t, cmd)
	assert.ErrorIs(t, err, errs.ErrValueIsRequired)
}

func TestUserRegistrationCommandHandler_Failure_EmptyTelegramChatID(t *testing.T) {
	cmd, err := commands.NewUserRegistrationCommand("test", "0", user.ProviderTelegram)
	assert.Nil(t, cmd)
	assert.ErrorIs(t, err, errs.ErrValueIsRequired)
}

func TestUserRegistrationCommandHandler_Idempotent(t *testing.T) {
	ctx, pool := setupTest(t)

	cr := cmd2.NewCompositionRoot(configs.Config{}, pool)

	uowFactory := cr.NewUnitOfWorkFactory()

	handler := cr.NewUserRegistrationCommandHandler()

	cmd, err := commands.NewUserRegistrationCommand(
		"test",
		"123",
		user.ProviderTelegram,
	)
	require.NoError(t, err)

	// первый вызов
	err = handler.Handle(ctx, cmd)
	require.NoError(t, err)

	// второй вызов (повторный /start)
	err = handler.Handle(ctx, cmd)
	require.NoError(t, err)

	uow, err := uowFactory.New()
	require.NoError(t, err)

	// пользователь всё ещё один
	user1, err := uow.UserRepository().
		FindByExternalProvider(user.ProviderTelegram, "123")
	require.NoError(t, err)
	require.NotNil(t, user1)

	usersCount := 0
	err = pool.GetContext(ctx, &usersCount, `SELECT COUNT(*) FROM users`)
	require.NoError(t, err)
	assert.Equal(t, 1, usersCount)

	identitiesCount := 0
	err = pool.GetContext(ctx, &identitiesCount, `SELECT COUNT(*) FROM external_identities`)
	require.NoError(t, err)
	assert.Equal(t, 1, identitiesCount)
}
