package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
	"google.golang.org/protobuf/proto"
)

func TestGetAccounts(t *testing.T) {
	ctx := context.TODO()
	db, err := dbHandle(ctx)
	assert.NoError(t, err)
	// create the schema utils
	schemaUtils := NewSchemaUtils(db)
	// create the accounts table
	require.NoError(t, schemaUtils.CreateAccountsTable(ctx))

	// let insert some accounts record into the database
	insertStatement := `
	INSERT INTO accounts(account_id, account_balance, account_owner, is_closed)
	VALUES 
	    ('account-1', 500.21, 'John Doe', TRUE),
	    ('account-2', 200.00, 'Mr Smith', FALSE),
	    ('account-3', 1000.00, 'Lady G.', FALSE),
	    ('account-4', 250.00, 'Mrs Peng', FALSE);
	`

	_, err = db.Exec(ctx, insertStatement)
	require.NoError(t, err)

	// create the storage for test
	storage := NewTestStorage(db)

	// define the accounts we want to fetch
	accountIDs := []string{"account-1", "account-5", "account-3"}
	accounts, err := storage.GetAccounts(ctx, accountIDs)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, 3)

	// let us define the expected
	account1 := &pb.BankAccount{
		AccountId:      "account-1",
		AccountBalance: 500.21,
		AccountOwner:   "John Doe",
		IsClosed:       true,
	}
	account3 := &pb.BankAccount{
		AccountId:      "account-3",
		AccountBalance: 1000.00,
		AccountOwner:   "Lady G.",
		IsClosed:       false,
	}

	expecteds := []*pb.BankAccount{
		account1,
		nil,
		account3,
	}

	for index, account := range accounts {
		require.True(t, proto.Equal(expecteds[index], account))
	}

	// free resources
	assert.NoError(t, schemaUtils.DropAccountsTable(ctx))
	assert.NoError(t, db.Disconnect(ctx))
}
