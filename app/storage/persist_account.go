package storage

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/tochemey/gopack/otel/trace"
	"github.com/tochemey/gopack/postgres"
	"google.golang.org/protobuf/proto"

	"github.com/tochemey/cos-go-sample/app/log"
	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
)

// PersistAccount persist an account record into the database
func (s *storage) PersistAccount(ctx context.Context, account *pb.BankAccount) error {
	// set the observability span
	spanCtx, span := trace.SpanContext(ctx, "PersistAccount")
	defer span.End()

	// get the context logger
	logger := log.WithContext(ctx)

	// check whether the account record is set or not
	if account == nil || proto.Equal(account, new(pb.BankAccount)) {
		err := errors.New("the account data record is not set")
		logger.Error(err)
		return err
	}

	// start a transaction runner
	txRunner, err := postgres.NewTxRunner(spanCtx, s.db)
	// handle the error
	if err != nil {
		logger.Error(errors.Wrap(err, "failed to setup database transaction"))
		return err
	}

	// build the transaction runner
	runner := txRunner.
		AddQueryBuilder(&deleteStmt{account}).
		AddQueryBuilder(&insertionStateStmt{account})

	// handle the error
	if err = runner.Execute(); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

type deleteStmt struct {
	account *pb.BankAccount
}

// BuildQuery build the SQL statement and arguments to run against the database
func (s deleteStmt) BuildQuery() (sqlStatement string, args []any, err error) {
	// build the actual SQL statement and params
	sqlStatement, args, err = sq.
		StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Delete("accounts").
		Where(sq.Eq{"account_id": s.account.GetAccountId()}).
		ToSql()
	return
}

type insertionStateStmt struct {
	account *pb.BankAccount
}

// BuildQuery build the SQL statement and arguments to run against the database
func (s insertionStateStmt) BuildQuery() (sqlStatement string, args []any, err error) {
	// build the actual SQL statement and params
	sqlStatement, args, err = sq.
		StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Insert("accounts").
		Columns(
			"account_id",
			"account_balance",
			"account_owner",
			"is_closed").
		Values(
			s.account.GetAccountId(),
			s.account.GetAccountBalance(),
			s.account.GetAccountOwner(),
			s.account.GetIsClosed(),
		).
		ToSql()
	return
}
