package cos

import (
	"context"
	"fmt"

	gopack "github.com/tochemey/gopack/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
	cospb "github.com/tochemey/cos-go-sample/gen/chief_of_state/v1"
)

// Client is used by both the service and the consumer implementation.
type Client interface {
	ProcessCommand(ctx context.Context, accountID string, command proto.Message) (*pb.BankAccount, *cospb.MetaData, error)
	GetState(ctx context.Context, accountID string) (*pb.BankAccount, *cospb.MetaData, error)
}

// client implements the Client interface
type client struct {
	remote cospb.ChiefOfStateServiceClient
}

var _ Client = &client{}

// NewClient creates a new instance of Client
func NewClient(cosHost string, cosPort int) (Client, error) {
	// get the grpc client connection to CoS
	conn, err := gopack.DefaultClientConn(fmt.Sprintf("%v:%v", cosHost, cosPort))
	// handle the error
	if err != nil {
		return nil, err
	}
	return &client{
		remote: cospb.NewChiefOfStateServiceClient(conn),
	}, nil
}

// ProcessCommand sends a command to COS and returns the resulting state and metadata
func (c client) ProcessCommand(ctx context.Context, accountID string, command proto.Message) (*pb.BankAccount, *cospb.MetaData, error) {
	// require a command
	if command == nil {
		return nil, nil, status.Error(codes.Internal, "command is missing")
	}

	// pack command into Any
	cmdAny, _ := anypb.New(command)

	// construct COS request
	request := &cospb.ProcessCommandRequest{
		EntityId: accountID,
		Command:  cmdAny,
	}

	// call COS get response
	response, err := c.remote.ProcessCommand(ctx, request)
	if err != nil {
		return nil, nil, err
	}

	// unpack the resulting state
	resultingState, err := UnmarshalState(response.GetState())
	if err != nil {
		return nil, nil, err
	}

	// return the company and the metadata
	return resultingState, response.GetMeta(), nil
}

// GetState retrieves the current  state of an entity and its metadata
func (c client) GetState(ctx context.Context, accountID string) (*pb.BankAccount, *cospb.MetaData, error) {
	// call CoS
	response, err := c.remote.GetState(ctx, &cospb.GetStateRequest{EntityId: accountID})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.NotFound {
				return nil, nil, nil
			}
		}

		return nil, nil, err
	}

	// handle nil response like a NOT_FOUND
	if response == nil {
		return nil, nil, nil
	}

	// unpack the resulting state
	resultingState, err := UnmarshalState(response.GetState())
	if err != nil {
		return nil, nil, err
	}

	// return
	return resultingState, response.GetMeta(), nil
}

// UnmarshalState unpacks the actual state from the proto any message
func UnmarshalState(any *anypb.Any) (*pb.BankAccount, error) {
	msg, err := any.UnmarshalNew()
	if err != nil {
		return nil, err
	}

	switch v := msg.(type) {
	case *pb.BankAccount:
		return v, nil
	case *emptypb.Empty:
		return nil, nil
	default:
		expected := proto.MessageName(new(pb.BankAccount))
		return nil, status.Errorf(codes.Internal, "expecting %s got %s", expected, any.GetTypeUrl())
	}
}
