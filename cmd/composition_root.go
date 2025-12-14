package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"

	"github.com/Nemizar/coin_tamer_bot/internal/core/application/eventshandler"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/user"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/ddd"

	"github.com/Nemizar/coin_tamer_bot/configs"
	"github.com/Nemizar/coin_tamer_bot/internal/adapters/out/postgres"
	"github.com/Nemizar/coin_tamer_bot/internal/adapters/out/sl"
	"github.com/Nemizar/coin_tamer_bot/internal/adapters/out/sl/handlers/slogpretty"
	"github.com/Nemizar/coin_tamer_bot/internal/core/application/usecases/commands"
	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"
)

type CompositionRoot struct {
	config  configs.Config
	db      *sqlx.DB
	logger  ports.Logger
	closers []Closer
}

func NewCompositionRoot(config configs.Config, db *sqlx.DB) *CompositionRoot {
	cr := &CompositionRoot{
		config:  config,
		db:      db,
		logger:  setupLogger(config),
		closers: make([]Closer, 0),
	}

	cr.RegisterCloser(db)

	return cr
}

func (cr *CompositionRoot) NewUnitOfWork() ports.UnitOfWork {
	unitOfWork, err := postgres.NewUnitOfWork(cr.db, cr.NewMediatrWithSubscriptions(), cr.Logger())
	if err != nil {
		panic(fmt.Sprintf("cannot create UnitOfWork: %v", err))
	}

	return unitOfWork
}

func (cr *CompositionRoot) NewUnitOfWorkFactory() ports.UnitOfWorkFactory {
	unitOfWorkFactory, err := postgres.NewUnitOfWorkFactory(cr.db, cr.NewMediatrWithSubscriptions(), cr.Logger())
	if err != nil {
		panic(fmt.Sprintf("cannot create UnitOfWorkFactory: %v", err))
	}

	return unitOfWorkFactory
}

func (cr *CompositionRoot) NewUserRegistrationCommandHandler() commands.UserRegistrationCommandHandler {
	handler, err := commands.NewUserRegistrationCommandHandler(cr.logger, cr.NewUnitOfWorkFactory())
	if err != nil {
		panic(fmt.Sprintf("cannot create UserRegistrationCommandHandler: %v", err))
	}

	return handler
}

func (cr *CompositionRoot) NewMediatrWithSubscriptions() ddd.Mediatr {
	mediatr := ddd.NewMediatr()
	mediatr.Subscribe(cr.NewUserCreatedDomainEventHandler(), user.NewEmptyCreateEvent())

	return mediatr
}

func (cr *CompositionRoot) NewUserCreatedDomainEventHandler() ddd.EventHandler {
	return eventshandler.NewUserCreatedEventHandler()
}

func (cr *CompositionRoot) Logger() ports.Logger {
	return cr.logger
}

func setupLogger(c configs.Config) ports.Logger {
	var log ports.Logger

	if c.IsProd() {
		log = sl.NewSlogLogger(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
	} else {
		log = setupPrettySlog()
	}

	return log
}

func setupPrettySlog() ports.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return sl.NewSlogLogger(slog.New(handler))
}
