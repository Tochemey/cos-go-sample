package commands

import (
	"context"

	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
	cospb "github.com/tochemey/cos-go-sample/gen/chief_of_state/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var (
	errCommandNotDefined        = status.Error(codes.Internal, "the command is not defined")
	errMissingPriorState        = status.Error(codes.InvalidArgument, "the priorState is not defined")
	errCommandSentToWrongEntity = status.Error(codes.InvalidArgument, "the command is sent to the wrong entity")
	errUnhandledCommand         = func(command proto.Message) error {
		return status.Errorf(codes.Internal, "received unhandled command (%s)", command.ProtoReflect().Descriptor().FullName())
	}
)

type Dispatcher interface {
	// Dispatch dispatches the given command and return the appropriate event or an error
	Dispatch(ctx context.Context, command proto.Message, priorState *pb.BankAccount, priorMeta *cospb.MetaData) (event proto.Message, err error)
}

type dispatcher struct{}

var _ Dispatcher = (*dispatcher)(nil)

// NewDispatcher create an instance of Dispatcher
func NewDispatcher() Dispatcher {
	return &dispatcher{}
}

// Dispatch dispatches the given command and return the appropriate event or an error
func (h dispatcher) Dispatch(ctx context.Context, command proto.Message, priorState *pb.BankAccount, priorMeta *cospb.MetaData) (event proto.Message, err error) { //nolint
	switch typedCmd := command.(type) {
	case *pb.OpenAccount:
		return openAccount(ctx, typedCmd)
	case *pb.CreditAccount:
		return creditAccount(ctx, typedCmd, priorState)
	case *pb.DebitAccount:
		return debitAccount(ctx, typedCmd, priorState)
	case nil:
		return nil, errCommandNotDefined
	default:
		return nil, errUnhandledCommand(typedCmd)
	}
}
