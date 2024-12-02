package commands

import (
	"context"

	"github.com/tochemey/gopack/otel/trace"
	"google.golang.org/protobuf/proto"

	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
)

// openAccount handles the Open Account command. When the command is valid the account credited event is returned
// to be persisted. On the contrary a validation error is returned
func openAccount(ctx context.Context, command *pb.OpenAccount) (*pb.AccountOpened, error) {
	// add a span context to trace the event handler
	ctx, span := trace.SpanContext(ctx, "HandleOpenAccount")
	defer span.End()

	// let us make a copy of the command
	commandCopy := proto.Clone(command).(*pb.OpenAccount)

	return &pb.AccountOpened{
		AccountId:    commandCopy.GetAccountId(),
		Balance:      commandCopy.GetOpeningBalance(),
		AccountOwner: commandCopy.GetAccountOwner(),
	}, nil
}
