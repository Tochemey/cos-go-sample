package cos

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
	cospb "github.com/tochemey/cos-go-sample/gen/chief_of_state/v1"
	mocks "github.com/tochemey/cos-go-sample/mocks/gen/chief_of_state/v1"
)

type cosClientTestSuite struct {
	suite.Suite
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestCosClient(t *testing.T) {
	suite.Run(t, new(cosClientTestSuite))
}

func (s *cosClientTestSuite) TestNewClient() {
	s.Run("happy path", func() {
		// create a context
		ctx := context.TODO()
		// this will work because grpc connection won't wait for connections to be
		// established, and connecting happens in the background
		cosClient, err := NewClient(ctx, "localhost", 50051)
		s.Assert().NotNil(cosClient)
		s.Assert().NoError(err)
	})
}

func (s *cosClientTestSuite) TestUnmarshalState() {
	s.Run("with valid state", func() {
		// create a new state
		state := new(pb.BankAccount)
		// pack that state into anypb
		anypbState, err := anypb.New(state)
		s.Assert().NoError(err)
		s.Assert().NotNil(anypbState)

		unpacked, err := UnmarshalState(anypbState)
		s.Assert().NoError(err)
		s.Assert().True(proto.Equal(state, unpacked))
	})
	s.Run("with an empty proto message", func() {
		// create an empty proto message
		empty := new(emptypb.Empty)
		// pack into anypb
		anypbState, err := anypb.New(empty)
		s.Assert().NoError(err)
		s.Assert().NotNil(anypbState)
		unpacked, err := UnmarshalState(anypbState)
		s.Assert().NoError(err)
		s.Assert().Nil(unpacked)
	})
	s.Run("with invalid state", func() {
		// create a wrong state
		anypbState, err := anypb.New(wrapperspb.String("not a valid state"))
		s.Assert().NoError(err)
		s.Assert().NotNil(anypbState)
		unpacked, err := UnmarshalState(anypbState)
		s.Assert().Error(err)
		s.Assert().Nil(unpacked)
	})
	s.Run("with invalid anypb state", func() {
		// create a wrong state
		anypbState := &anypb.Any{
			TypeUrl: "",
			Value:   nil,
		}
		unpacked, err := UnmarshalState(anypbState)
		s.Assert().Error(err)
		s.Assert().Nil(unpacked)
	})
}

func (s *cosClientTestSuite) TestProcessCommand() {
	s.Run("with nil command", func() {
		// create the remote client
		mockRemoteClient := &mocks.ChiefOfStateServiceClient{}
		mockCos := client{remote: mockRemoteClient}
		state, meta, err := mockCos.ProcessCommand(context.TODO(), uuid.NewString(), nil)
		expectedError := status.Error(codes.Internal, "command is missing")
		s.Assert().Nil(state)
		s.Assert().Nil(meta)
		s.Assert().EqualError(err, expectedError.Error())
	})
	s.Run("with happy path", func() {
		ctx := context.TODO()
		accountID := "account-1"
		accountBal := 150.55
		accountOwner := "John Doe"
		amount := 50.00

		// create the prior state
		currentState := &pb.BankAccount{
			AccountId:      accountID,
			AccountBalance: accountBal,
			AccountOwner:   accountOwner,
			IsClosed:       false,
		}

		now := timestamppb.Now()
		anypbState, err := anypb.New(currentState)
		s.Assert().NoError(err)
		cosMeta := &cospb.MetaData{
			EntityId:       accountID,
			RevisionNumber: 2,
			RevisionDate:   now,
		}
		// create the process command response
		cosResp := &cospb.ProcessCommandResponse{State: anypbState, Meta: cosMeta}
		// create the remote client
		mockRemoteClient := &mocks.ChiefOfStateServiceClient{}
		mockRemoteClient.On("ProcessCommand", ctx, mock.Anything).Return(cosResp, nil)
		// create the CoS client
		mockCos := client{mockRemoteClient}
		// create the command
		cmd := &pb.CreditAccount{
			AccountId: accountID,
			Amount:    amount,
		}
		state, meta, err := mockCos.ProcessCommand(ctx, accountID, cmd)
		s.Assert().NoError(err)
		s.Assert().NotNil(meta)
		s.Assert().NotNil(state)
		s.Assert().True(proto.Equal(currentState, state))
		s.Assert().True(proto.Equal(cosMeta, meta))
		mockRemoteClient.AssertExpectations(s.T())
	})
	s.Run("with remote client failure", func() {
		ctx := context.TODO()
		accountID := "account-1"
		amount := 50.00

		// create the remote client
		mockRemoteClient := &mocks.ChiefOfStateServiceClient{}
		mockRemoteClient.On("ProcessCommand", ctx, mock.Anything).Return(nil, status.Error(codes.Internal, ""))
		// create the CoS client
		mockCos := client{mockRemoteClient}
		cmd := &pb.CreditAccount{
			AccountId: accountID,
			Amount:    amount,
		}

		state, meta, err := mockCos.ProcessCommand(ctx, accountID, cmd)
		s.Assert().Error(err)
		s.Assert().Nil(meta)
		s.Assert().Nil(state)
		mockRemoteClient.AssertExpectations(s.T())
	})
	s.Run("with invalid state returned", func() {
		ctx := context.TODO()
		// create the various ID
		accountID := "account-1"
		amount := 50.00
		now := timestamppb.Now()

		anypbState, err := anypb.New(wrapperspb.String("not a valid state"))
		s.Assert().NoError(err)
		cosMeta := &cospb.MetaData{
			EntityId:       accountID,
			RevisionNumber: 2,
			RevisionDate:   now,
		}
		// create the process command response
		cosResp := &cospb.ProcessCommandResponse{State: anypbState, Meta: cosMeta}
		// create the remote client
		mockRemoteClient := &mocks.ChiefOfStateServiceClient{}
		mockRemoteClient.On("ProcessCommand", ctx, mock.Anything).Return(cosResp, nil)
		// create the CoS client
		mockCos := client{mockRemoteClient}
		cmd := &pb.CreditAccount{
			AccountId: accountID,
			Amount:    amount,
		}
		state, meta, err := mockCos.ProcessCommand(ctx, accountID, cmd)
		s.Assert().Error(err)
		s.Assert().Nil(meta)
		s.Assert().Nil(state)
		mockRemoteClient.AssertExpectations(s.T())
	})
}

func (s *cosClientTestSuite) TestGetState() {
	s.Run("with happy path", func() {
		ctx := context.TODO()
		accountID := uuid.NewString()
		now := timestamppb.Now()
		// create the current state
		currentState := &pb.BankAccount{AccountId: accountID}
		anypbState, err := anypb.New(currentState)
		s.Assert().NoError(err)
		s.Assert().NotNil(anypbState)
		cosMeta := &cospb.MetaData{
			EntityId:       accountID,
			RevisionNumber: 2,
			RevisionDate:   now,
		}
		// create the process command response
		cosResp := &cospb.GetStateResponse{State: anypbState, Meta: cosMeta}
		// create the remote client
		mockRemoteClient := &mocks.ChiefOfStateServiceClient{}
		mockRemoteClient.On("GetState", ctx, mock.Anything).Return(cosResp, nil)
		// create the CoS client
		mockCos := client{remote: mockRemoteClient}
		state, meta, err := mockCos.GetState(ctx, accountID)
		s.Assert().NoError(err)
		s.Assert().NotNil(meta)
		s.Assert().NotNil(state)
		s.Assert().True(proto.Equal(currentState, state))
		s.Assert().True(proto.Equal(cosMeta, meta))
		mockRemoteClient.AssertExpectations(s.T())
	})
	s.Run("with CoS failure", func() {
		ctx := context.TODO()
		accountID := uuid.NewString()

		// create the current state
		currentState := &pb.BankAccount{AccountId: accountID}
		anypbState, err := anypb.New(currentState)
		s.Assert().NoError(err)
		s.Assert().NotNil(anypbState)

		// create the remote client
		mockRemoteClient := &mocks.ChiefOfStateServiceClient{}
		mockRemoteClient.On("GetState", ctx, mock.Anything).Return(nil, status.Error(codes.Unavailable, ""))
		// create the CoS client
		mockCos := client{remote: mockRemoteClient}
		state, meta, err := mockCos.GetState(ctx, accountID)
		s.Assert().Error(err)
		s.Assert().Nil(meta)
		s.Assert().Nil(state)
		mockRemoteClient.AssertExpectations(s.T())
	})
	s.Run("with invalid state", func() {
		ctx := context.TODO()
		accountID := uuid.NewString()
		now := timestamppb.Now()
		// create the current state
		anypbState, err := anypb.New(wrapperspb.String("not a valid state"))
		s.Assert().NoError(err)
		s.Assert().NotNil(anypbState)
		s.Assert().NoError(err)
		s.Assert().NotNil(anypbState)
		cosMeta := &cospb.MetaData{
			EntityId:       accountID,
			RevisionNumber: 2,
			RevisionDate:   now,
		}
		// create the process command response
		cosResp := &cospb.GetStateResponse{State: anypbState, Meta: cosMeta}
		// create the remote client
		mockRemoteClient := &mocks.ChiefOfStateServiceClient{}
		mockRemoteClient.On("GetState", ctx, mock.Anything).Return(cosResp, nil)
		// create the CoS client
		mockCos := client{remote: mockRemoteClient}
		state, meta, err := mockCos.GetState(ctx, accountID)
		s.Assert().Error(err)
		s.Assert().Nil(meta)
		s.Assert().Nil(state)
		mockRemoteClient.AssertExpectations(s.T())
	})
	s.Run("with not found", func() {
		ctx := context.TODO()
		accountID := uuid.NewString()

		// create the remote client
		mockRemoteClient := &mocks.ChiefOfStateServiceClient{}
		mockRemoteClient.On("GetState", ctx, mock.Anything).Return(nil, status.Error(codes.NotFound, "state not found"))
		// create the CoS client
		mockCos := client{remote: mockRemoteClient}
		state, meta, err := mockCos.GetState(ctx, accountID)
		s.Assert().NoError(err)
		s.Assert().Nil(meta)
		s.Assert().Nil(state)
		mockRemoteClient.AssertExpectations(s.T())
	})
	s.Run("with nil response", func() {
		ctx := context.TODO()
		accountID := uuid.NewString()
		// create the remote client
		mockRemoteClient := &mocks.ChiefOfStateServiceClient{}
		mockRemoteClient.On("GetState", ctx, mock.Anything).Return(nil, nil)
		// create the CoS client
		mockCos := client{remote: mockRemoteClient}
		state, meta, err := mockCos.GetState(ctx, accountID)
		s.Assert().NoError(err)
		s.Assert().Nil(meta)
		s.Assert().Nil(state)
		mockRemoteClient.AssertExpectations(s.T())
	})
}
