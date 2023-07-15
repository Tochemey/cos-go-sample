package storage

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/tochemey/gopack/log/zapl"
	"github.com/tochemey/gopack/postgres"
)

type storage struct {
	db postgres.Postgres
	sb sq.StatementBuilderType
}

// enforce compilation error
var _ Storage = (*storage)(nil)

// New creates an instance of Storage
// This call will panic when there is an error which expected since this call will be done on application start
func New(ctx context.Context) Storage {
	// load the database configuration from environment variables
	config := LoadConfig()
	// create the database connection
	db := postgres.New(config)
	// connect to the database
	if err := db.Connect(ctx); err != nil {
		zapl.Panic(errors.Wrap(err, "failed to connect to the postgres database"))
	}
	// create the instance and return it
	return &storage{
		db: db,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// NewTestStorage creates an instance of Storage purposely for unit tests or
// integration tests. Never use this in production code
func NewTestStorage(db postgres.Postgres) Storage {
	return &storage{
		db: db,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Shutdown shuts down the database connection
func (s *storage) Shutdown(ctx context.Context) func(ctx context.Context) error {
	// prepare the function to run
	return func(ctx context.Context) error {
		if err := s.db.Disconnect(ctx); err != nil {
			return err
		}
		return nil
	}
}
