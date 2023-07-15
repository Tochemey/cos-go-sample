package storage

import (
	"context"

	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
)

// Storage represents the storage API
type Storage interface {
	Shutdown(ctx context.Context) func()
	PersistAccount(ctx context.Context, account *pb.BankAccount) error
}
