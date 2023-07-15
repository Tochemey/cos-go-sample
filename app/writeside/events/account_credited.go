package events

import (
	"context"

	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
	"github.com/tochemey/gopack/otel/trace"
	"google.golang.org/protobuf/proto"
)

// accountCredited handles the account credited event and return the resulting state
func accountCredited(ctx context.Context, event *pb.AccountCredited, priorState *pb.BankAccount) (*pb.BankAccount, error) {
	// add a span context to trace the event handler
	ctx, span := trace.SpanContext(ctx, "HandleAccountCredited")
	defer span.End()

	eventCopy := proto.Clone(event).(*pb.AccountCredited)
	stateCopy := proto.Clone(priorState).(*pb.BankAccount)

	stateCopy.AccountBalance = stateCopy.GetAccountBalance() + eventCopy.GetAmount()

	return stateCopy, nil
}
