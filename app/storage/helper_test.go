package storage

import (
	"context"
	"os"
	"testing"

	"github.com/tochemey/gopack/postgres"
)

var testContainer *postgres.TestContainer

const (
	testUser             = "test"
	testDatabase         = "testdb"
	testDatabasePassword = "test"
)

// TestMain will spawn a postgres database container that will be used for all tests
// making use of the postgres database container
func TestMain(m *testing.M) {
	// set the test container
	testContainer = postgres.NewTestContainer(testDatabase, testUser, testDatabasePassword)
	// execute the tests
	code := m.Run()
	// free resources
	testContainer.Cleanup()
	// exit the tests
	os.Exit(code)
}

// dbHandle returns a test db
func dbHandle(ctx context.Context) (*postgres.TestDB, error) {
	db := testContainer.GetTestDB()
	if err := db.Connect(ctx); err != nil {
		return nil, err
	}
	return db, nil
}

// SchemaUtils help create the various test tables in unit/integration tests
type SchemaUtils struct {
	db *postgres.TestDB
}

// NewSchemaUtils creates an instance of SchemaUtils
func NewSchemaUtils(db *postgres.TestDB) *SchemaUtils {
	return &SchemaUtils{db: db}
}

// CreateAccountsTable creates the accounts table used for unit and integration tests
func (s SchemaUtils) CreateAccountsTable(ctx context.Context) error {
	schemaDDL := `
	DROP TABLE IF EXISTS accounts;
	-- accounts relation
	CREATE TABLE accounts(
		account_id VARCHAR(255) NOT NULL,
		account_balance NUMERIC(19, 2) NOT NULL,
		account_owner VARCHAR(255) NOT NULL,
		is_closed BOOLEAN NOT NULL,
	
		PRIMARY KEY (account_id)
	);
	`
	_, err := s.db.Exec(ctx, schemaDDL)
	return err
}

// DropAccountsTable drops the accounts table used in unit test
// This is useful for resource cleanup after a unit test
func (s SchemaUtils) DropAccountsTable(ctx context.Context) error {
	return s.db.DropTable(ctx, "accounts")
}
