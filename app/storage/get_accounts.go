package storage

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/tochemey/gopack/otel/trace"

	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
)

// GetAccounts fetches the ordered list of accounts.
// The order in the list is the same as the order of the account's ids sent.
// When a record is not found nil is return instead
func (s *storage) GetAccounts(ctx context.Context, accountIDs []string) (accounts []*pb.BankAccount, err error) {
	// set the observability span
	spanCtx, span := trace.SpanContext(ctx, "GetAccounts")
	defer span.End()

	// create the select statement
	statement := s.sb.
		Select(
			"account_id",
			"account_balance",
			"account_owner",
			"is_closed").
		From("accounts").
		Where(sq.Eq{"account_id": accountIDs})

	// get the sql statement and the arguments
	query, args, err := statement.ToSql()
	// handle the error
	if err != nil {
		return nil, errors.Wrap(err, "unable to build sql statement")
	}

	// define the data type to hold the records fetched from the database
	type row struct {
		AccountID      string
		AccountBalance float64
		AccountOwner   string
		IsClosed       bool
	}

	// create the variable to hold the scanned account records
	var rows []*row
	// fetch the data and handle the eventual select error
	if err = s.db.SelectAll(spanCtx, &rows, query, args...); err != nil {
		return nil, errors.Wrap(err, "failed to fetch account records")
	}

	// create a map holding account id and record scanned from the database
	recordsMap := make(map[string]*pb.BankAccount)
	// iterate the rows scanned and build the map
	for _, row := range rows {
		recordsMap[row.AccountID] = &pb.BankAccount{
			AccountId:      row.AccountID,
			AccountBalance: row.AccountBalance,
			AccountOwner:   row.AccountOwner,
			IsClosed:       row.IsClosed,
		}
	}

	// initialize the output data
	accounts = make([]*pb.BankAccount, len(accountIDs))
	// set the output data with the records fetched
	for orderNr, accountID := range accountIDs {
		// look up for the account id in the record map and sets
		// the corresponding record fetched.
		if account, found := recordsMap[accountID]; found {
			accounts[orderNr] = account
		}
	}

	return
}
