package commands

import (
	"context"

	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
	"github.com/tochemey/gopack/log/zapl"
	"github.com/tochemey/gopack/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// debitAccount handles the Debit Account command. When the command is valid the account debited event is returned
// to be persisted. On the contrary a validation error is returned
func debitAccount(ctx context.Context, command *pb.DebitAccount, priorState *pb.BankAccount) (*pb.AccountDebited, error) {
	// add a span context to trace the event handler
	ctx, span := trace.SpanContext(ctx, "HandleDebitAccount")
	defer span.End()

	// get the context logger
	logger := zapl.WithContext(ctx)

	// let us make a copy of the command and the prior state
	commandCopy := proto.Clone(command).(*pb.DebitAccount)
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

	// perform some validation
	balanceAfter := priorStateCopy.GetAccountBalance() - commandCopy.GetAmount()
	// return a validation error when the balance after is negative or zero
	if balanceAfter <= 0 {
		logger.Warn("insufficient balance")
		return nil, status.Error(codes.InvalidArgument, "insufficient balance")
	}

	// create the account debited event to persist into the data store
	return &pb.AccountDebited{
		AccountId: commandCopy.GetAccountId(),
		Amount:    commandCopy.GetAmount(),
	}, nil
}
