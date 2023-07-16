package dbwriter

import (
	"context"

	"github.com/pkg/errors"
	"github.com/tochemey/cos-go-sample/app/cos"
	"github.com/tochemey/cos-go-sample/app/log"
	"github.com/tochemey/cos-go-sample/app/storage"
	cospb "github.com/tochemey/cos-go-sample/gen/chief_of_state/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Service is an implementation of the CoS ReadSide handler interface
type Service struct {
	dataStore storage.Storage
}

// NewService creates a new instance of service
func NewService(dataStore storage.Storage) (*Service, error) {
	// check whether the data store is defined or not
	if dataStore == nil {
		return nil, errors.New("the dataStore is not defined")
	}

	// return the new instance of Service
	return &Service{
		dataStore: dataStore,
	}, nil
}

// RegisterService registers the gRPC service
func (s Service) RegisterService(sv *grpc.Server) {
	cospb.RegisterReadSideHandlerServiceServer(sv, s)
}

// HandleReadSide handles read-side requests
func (s Service) HandleReadSide(ctx context.Context, request *cospb.HandleReadSideRequest) (*cospb.HandleReadSideResponse, error) {
	// set the logger with the context
	logger := log.WithContext(ctx)
	// make a copy of the request
	requestCopy := proto.Clone(request).(*cospb.HandleReadSideRequest)
	// check whether the user is nil
	if requestCopy.GetState() == nil {
		err := errors.New("the account state is not set")
		logger.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// let us unmarshall the user
	unpackState, err := cos.UnmarshalState(request.GetState())
	// handle the error
	if err != nil {
		err = errors.Wrapf(err, "failed to unpack state:(%s)", request.GetState().GetTypeUrl())
		logger.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// persist the data into the data store
	if err = s.dataStore.PersistAccount(ctx, unpackState); err != nil {
		err := errors.Wrap(err, "failed to persist account into the data store")
		logger.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	// return the successful handling of the read-side request
	return &cospb.HandleReadSideResponse{Successful: true}, nil
}
