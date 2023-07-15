package commands

import (
	"context"

	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
	"github.com/tochemey/gopack/log/zapl"
	"github.com/tochemey/gopack/otel/trace"
	"google.golang.org/protobuf/proto"
)

// creditAccount handles the Credit Account command. When the command is valid the account credited event is returned
// to be persisted. On the contrary a validation error is returned
func creditAccount(ctx context.Context, command *pb.CreditAccount, priorState *pb.BankAccount) (*pb.AccountCredited, error) {
	// add a span context to trace the event handler
	ctx, span := trace.SpanContext(ctx, "HandleCreditAccount")
	defer span.End()

	// get the context logger
	logger := zapl.WithContext(ctx)

	// let us make a copy of the command and the prior state
	commandCopy := proto.Clone(command).(*pb.CreditAccount)
	priorStateCopy := proto.Clone(priorState).(*pb.BankAccount)

	// check whether the prior state is defined or not
	if priorStateCopy == nil || proto.Equal(priorStateCopy, new(pb.BankAccount)) {
		// log a message for debugging purpose
		logger.Error("the prior state is not defined")
		// return an error of missing prior state
		return nil, errMissingPriorState
	}

	// let us verify that the command is sent to right aggregate.
	// this scenario with never occur but sanity check requires such verification
	if command.GetAccountId() != priorState.GetAccountId() {
		logger.Errorf("the account state:(%s) is not found", command.GetAccountId())
		return nil, errCommandSentToWrongEntity
	}

	// create the account credited event to persist into the data store
	return &pb.AccountCredited{
		AccountId: commandCopy.GetAccountId(),
		Amount:    commandCopy.GetAmount(),
	}, nil
}
