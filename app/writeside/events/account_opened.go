package events

import (
	"context"

	"github.com/tochemey/gopack/otel/trace"
	"google.golang.org/protobuf/proto"

	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
)

// accountOpened handles the account opened event and returns the resulting state
func accountOpened(ctx context.Context, event *pb.AccountOpened) (*pb.BankAccount, error) {
	// add a span context to trace the event handler
	ctx, span := trace.SpanContext(ctx, "HandleAccountOpened")
	defer span.End()

	// let us make a copy of the event
	eventCopy := proto.Clone(event).(*pb.AccountOpened)

	// return the resulting state
	return &pb.BankAccount{
		AccountId:      eventCopy.GetAccountId(),
		AccountBalance: eventCopy.GetBalance(),
		AccountOwner:   eventCopy.GetAccountOwner(),
		IsClosed:       false,
	}, nil
}
