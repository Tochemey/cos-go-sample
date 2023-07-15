package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
	"google.golang.org/protobuf/proto"
)

func TestPersistAccount(t *testing.T) {
	t.Run("With valid account record", func(t *testing.T) {
		ctx := context.TODO()
		db, err := dbHandle(ctx)
		assert.NoError(t, err)
		// create the schema utils
		schemaUtils := NewSchemaUtils(db)
		// create the accounts table
		require.NoError(t, schemaUtils.CreateAccountsTable(ctx))

		// create the storage for test
		storage := NewTestStorage(db)

		// create the account record to persist
		accountID := "account-1"
		accountBal := 150.55
		accountOwner := "John Doe"
		account := &pb.BankAccount{
			AccountId:      accountID,
			AccountBalance: accountBal,
			AccountOwner:   accountOwner,
			IsClosed:       false,
		}

		// persist the account
		require.NoError(t, storage.PersistAccount(ctx, account))

		// fetch the record
		accounts, err := storage.GetAccounts(ctx, []string{accountID})
		require.NoError(t, err)
		require.NotEmpty(t, accounts)
		require.Len(t, accounts, 1)

		assert.True(t, proto.Equal(account, accounts[0]))

		// free resources
		assert.NoError(t, schemaUtils.DropAccountsTable(ctx))
		assert.NoError(t, db.Disconnect(ctx))
	})
	t.Run("With invalid account record", func(t *testing.T) {
		ctx := context.TODO()
		db, err := dbHandle(ctx)
		assert.NoError(t, err)
		// create the schema utils
		schemaUtils := NewSchemaUtils(db)
		// create the accounts table
		require.NoError(t, schemaUtils.CreateAccountsTable(ctx))

		// create the storage for test
		storage := NewTestStorage(db)

		// create the account record to persist
		account := new(pb.BankAccount)
		// persist the account
		require.Error(t, storage.PersistAccount(ctx, account))

		// free resources
		assert.NoError(t, schemaUtils.DropAccountsTable(ctx))
		assert.NoError(t, db.Disconnect(ctx))
	})
}
