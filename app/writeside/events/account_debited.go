package events

import (
	"context"

	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
	"github.com/tochemey/gopack/otel/trace"
	"google.golang.org/protobuf/proto"
)

// accountDebited handles the account debited event and return the resulting state
func accountDebited(ctx context.Context, event *pb.AccountDebited, priorState *pb.BankAccount) (*pb.BankAccount, error) {
	// add a span context to trace the event handler
	ctx, span := trace.SpanContext(ctx, "HandleAccountDebited")
	defer span.End()

	eventCopy := proto.Clone(event).(*pb.AccountDebited)
	stateCopy := proto.Clone(priorState).(*pb.BankAccount)

	bal := stateCopy.GetAccountBalance() - eventCopy.GetAmount()
	stateCopy.AccountBalance = bal

	return stateCopy, nil
}
