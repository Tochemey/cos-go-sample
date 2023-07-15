package events

import (
	"context"

	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
	cospb "github.com/tochemey/cos-go-sample/gen/chief_of_state/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var (
	errEventNotDefined = status.Error(codes.Internal, "the event is not defined")
	errUnhandledEvent  = func(event proto.Message) error {
		return status.Errorf(codes.Internal, "received unhandled command (%s)", event.ProtoReflect().Descriptor().FullName())
	}
)

type Dispatcher interface {
	// Dispatch dispatches the given event and return the appropriate event or an error
	Dispatch(ctx context.Context, event proto.Message, priorState *pb.BankAccount, eventMeta *cospb.MetaData) (newState *pb.BankAccount, err error)
}

type dispatcher struct{}

var _ Dispatcher = (*dispatcher)(nil)

// NewDispatcher create an instance of Dispatcher
func NewDispatcher() Dispatcher {
	return &dispatcher{}
}

// Dispatch dispatches the given event and return the appropriate event or an error
func (h dispatcher) Dispatch(ctx context.Context, event proto.Message, priorState *pb.BankAccount, eventMeta *cospb.MetaData) (newState *pb.BankAccount, err error) {
	switch typedEvent := event.(type) {
	case *pb.AccountOpened:
		return accountOpened(ctx, typedEvent, priorState)
	case *pb.AccountCredited:
		return accountCredited(ctx, typedEvent, priorState)
	case *pb.AccountDebited:
		return accountDebited(ctx, typedEvent, priorState)
	case nil:
		return nil, errEventNotDefined
	default:
		return nil, errUnhandledEvent(typedEvent)
	}
}
