package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gopack "github.com/tochemey/gopack/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
	cospb "github.com/tochemey/cos-go-sample/gen/chief_of_state/v1"
	mocks "github.com/tochemey/cos-go-sample/mocks/app/cos"
)

func TestService(t *testing.T) {
	t.Run("With new instance", func(t *testing.T) {
		cosClient := new(mocks.Client)
		svc := NewService(cosClient)
		assert.NotNil(t, svc)
		// assert the type of svc
		assert.IsType(t, &Service{}, svc)

		// assert that svc implement pb.BankAccountServiceServer
		var p interface{} = svc
		_, ok := p.(pb.BankAccountServiceServer)
		assert.True(t, ok)
	})
	t.Run("With service registration", func(t *testing.T) {
		fn := func() {
			svc := &Service{}
			// get in process grpc test server
			testServer, _ := gopack.TestServer(nil)
			svc.RegisterService(testServer)
		}

		assert.NotPanics(t, fn)
	})
	t.Run("With OpenAccount request", func(t *testing.T) {
		ctx := context.TODO()
		accountID := uuid.NewString()
		accountOwner := "Mr Account"
		openingBalance := 50.0

		// create the rpc request
		rpcReq := &pb.OpenAccountRequest{
			AccountOwner: accountOwner,
			Balance:      openingBalance,
			AccountId:    &accountID,
		}

		// create the command sent to the cos mock service
		command := &pb.OpenAccount{
			AccountId:      accountID,
			AccountOwner:   accountOwner,
			OpeningBalance: openingBalance,
		}

		// create the resulting state when cos finishes processing the command
		state := &pb.BankAccount{
			AccountId:      accountID,
			AccountBalance: openingBalance,
			AccountOwner:   accountOwner,
			IsClosed:       false,
		}

		// create the cos meta
		cosMeta := &cospb.MetaData{EntityId: accountID, RevisionNumber: 1, RevisionDate: timestamppb.Now()}

		// create the expected response
		expected := &pb.OpenAccountResponse{Account: state}

		// create a mock cos client
		cosClient := new(mocks.Client)
		cosClient.On("ProcessCommand", ctx, accountID, command).Return(state, cosMeta, nil)
		svc := NewService(cosClient)
		require.NotNil(t, svc)

		// process the request
		actual, err := svc.OpenAccount(ctx, rpcReq)
		require.NoError(t, err)
		require.NotNil(t, actual)
		assert.True(t, proto.Equal(expected, actual))
		cosClient.AssertExpectations(t)
	})
	t.Run("With OpenAccount request with cos client failure", func(t *testing.T) {
		ctx := context.TODO()
		accountID := uuid.NewString()
		accountOwner := "Mr Account"
		openingBalance := 50.0

		// create the rpc request
		rpcReq := &pb.OpenAccountRequest{
			AccountOwner: accountOwner,
			Balance:      openingBalance,
			AccountId:    &accountID,
		}

		// create the command sent to the cos mock service
		command := &pb.OpenAccount{
			AccountId:      accountID,
			AccountOwner:   accountOwner,
			OpeningBalance: openingBalance,
		}

		// create the cos meta
		cosMeta := &cospb.MetaData{EntityId: accountID, RevisionNumber: 1, RevisionDate: timestamppb.Now()}

		// create the expected error
		err := status.Error(codes.DeadlineExceeded, "context canceled")

		// create a mock cos client
		cosClient := new(mocks.Client)
		cosClient.On("ProcessCommand", ctx, accountID, command).Return(new(pb.BankAccount), cosMeta, err)
		svc := NewService(cosClient)
		require.NotNil(t, svc)

		// process the request
		actual, err := svc.OpenAccount(ctx, rpcReq)
		require.Error(t, err)
		require.Nil(t, actual)
		cosClient.AssertExpectations(t)
	})
	t.Run("With DebitAccount request", func(t *testing.T) {
		ctx := context.TODO()
		accountID := uuid.NewString()
		accountOwner := "Mr Account"
		amount := 50.0

		// create the rpc request
		rpcReq := &pb.DebitAccountRequest{
			AccountId: accountID,
			Amount:    amount,
		}

		// create the command sent to the cos mock service
		command := &pb.DebitAccount{
			AccountId: accountID,
			Amount:    amount,
		}

		// create the resulting state when cos finishes processing the command
		state := &pb.BankAccount{
			AccountId:      accountID,
			AccountBalance: amount + 10.0,
			AccountOwner:   accountOwner,
			IsClosed:       false,
		}

		// create the cos meta
		cosMeta := &cospb.MetaData{EntityId: accountID, RevisionNumber: 1, RevisionDate: timestamppb.Now()}

		// create the expected response
		expected := &pb.DebitAccountResponse{Account: state}

		// create a mock cos client
		cosClient := new(mocks.Client)
		cosClient.On("ProcessCommand", ctx, accountID, command).Return(state, cosMeta, nil)
		svc := NewService(cosClient)
		require.NotNil(t, svc)

		// process the request
		actual, err := svc.DebitAccount(ctx, rpcReq)
		require.NoError(t, err)
		require.NotNil(t, actual)
		assert.True(t, proto.Equal(expected, actual))
		cosClient.AssertExpectations(t)
	})
	t.Run("With DebitAccount request with cos client failure", func(t *testing.T) {
		ctx := context.TODO()
		accountID := uuid.NewString()
		amount := 50.0

		// create the rpc request
		rpcReq := &pb.DebitAccountRequest{
			AccountId: accountID,
			Amount:    amount,
		}

		// create the command sent to the cos mock service
		command := &pb.DebitAccount{
			AccountId: accountID,
			Amount:    amount,
		}

		// create the cos meta
		cosMeta := &cospb.MetaData{EntityId: accountID, RevisionNumber: 1, RevisionDate: timestamppb.Now()}

		// create the expected response
		expectedErr := status.Error(codes.Unavailable, "service unavailable")

		// create a mock cos client
		cosClient := new(mocks.Client)
		cosClient.On("ProcessCommand", ctx, accountID, command).Return(new(pb.BankAccount), cosMeta, expectedErr)
		svc := NewService(cosClient)
		require.NotNil(t, svc)

		// process the request
		actual, err := svc.DebitAccount(ctx, rpcReq)
		require.Error(t, err)
		require.Nil(t, actual)
		assert.EqualError(t, err, expectedErr.Error())
		cosClient.AssertExpectations(t)
	})
	t.Run("With CreditAccount request", func(t *testing.T) {
		ctx := context.TODO()
		accountID := uuid.NewString()
		accountOwner := "Mr Account"
		amount := 50.0

		// create the rpc request
		rpcReq := &pb.CreditAccountRequest{
			AccountId: accountID,
			Amount:    amount,
		}

		// create the command sent to the cos mock service
		command := &pb.CreditAccount{
			AccountId: accountID,
			Amount:    amount,
		}

		// create the resulting state when cos finishes processing the command
		state := &pb.BankAccount{
			AccountId:      accountID,
			AccountBalance: amount + 10.0,
			AccountOwner:   accountOwner,
			IsClosed:       false,
		}

		// create the cos meta
		cosMeta := &cospb.MetaData{EntityId: accountID, RevisionNumber: 1, RevisionDate: timestamppb.Now()}

		// create the expected response
		expected := &pb.CreditAccountResponse{Account: state}

		// create a mock cos client
		cosClient := new(mocks.Client)
		cosClient.On("ProcessCommand", ctx, accountID, command).Return(state, cosMeta, nil)
		svc := NewService(cosClient)
		require.NotNil(t, svc)

		// process the request
		actual, err := svc.CreditAccount(ctx, rpcReq)
		require.NoError(t, err)
		require.NotNil(t, actual)
		assert.True(t, proto.Equal(expected, actual))
		cosClient.AssertExpectations(t)
	})
	t.Run("With CreditAccount request with cos client failure", func(t *testing.T) {
		ctx := context.TODO()
		accountID := uuid.NewString()
		amount := 50.0

		// create the rpc request
		rpcReq := &pb.CreditAccountRequest{
			AccountId: accountID,
			Amount:    amount,
		}

		// create the command sent to the cos mock service
		command := &pb.CreditAccount{
			AccountId: accountID,
			Amount:    amount,
		}

		// create the cos meta
		cosMeta := &cospb.MetaData{EntityId: accountID, RevisionNumber: 1, RevisionDate: timestamppb.Now()}

		// create the expected response
		expectedErr := status.Error(codes.Unavailable, "service unavailable")

		// create a mock cos client
		cosClient := new(mocks.Client)
		cosClient.On("ProcessCommand", ctx, accountID, command).Return(new(pb.BankAccount), cosMeta, expectedErr)
		svc := NewService(cosClient)
		require.NotNil(t, svc)

		// process the request
		actual, err := svc.CreditAccount(ctx, rpcReq)
		require.Error(t, err)
		require.Nil(t, actual)
		assert.EqualError(t, err, expectedErr.Error())
		cosClient.AssertExpectations(t)
	})
	t.Run("With GetAccount request", func(t *testing.T) {
		ctx := context.TODO()
		accountID := uuid.NewString()
		accountOwner := "Mr Account"
		openingBalance := 50.0

		// create the rpc request
		rpcReq := &pb.GetAccountRequest{
			AccountId: accountID,
		}

		// create the resulting state when cos finishes processing the command
		state := &pb.BankAccount{
			AccountId:      accountID,
			AccountBalance: openingBalance,
			AccountOwner:   accountOwner,
			IsClosed:       false,
		}

		// create the cos meta
		cosMeta := &cospb.MetaData{EntityId: accountID, RevisionNumber: 1, RevisionDate: timestamppb.Now()}

		// create the expected response
		expected := &pb.GetAccountResponse{Account: state}

		// create a mock cos client
		cosClient := new(mocks.Client)
		cosClient.On("GetState", ctx, accountID).Return(state, cosMeta, nil)
		svc := NewService(cosClient)
		require.NotNil(t, svc)

		// process the request
		actual, err := svc.GetAccount(ctx, rpcReq)
		require.NoError(t, err)
		require.NotNil(t, actual)
		assert.True(t, proto.Equal(expected, actual))
		cosClient.AssertExpectations(t)
	})
	t.Run("With GetAccount request with cos client failure", func(t *testing.T) {
		ctx := context.TODO()
		accountID := uuid.NewString()

		// create the rpc request
		rpcReq := &pb.GetAccountRequest{
			AccountId: accountID,
		}

		// create the cos meta
		cosMeta := &cospb.MetaData{EntityId: accountID, RevisionNumber: 1, RevisionDate: timestamppb.Now()}

		// create the expected error
		err := status.Error(codes.DeadlineExceeded, "context canceled")

		// create a mock cos client
		cosClient := new(mocks.Client)
		cosClient.On("GetState", ctx, accountID).Return(new(pb.BankAccount), cosMeta, err)
		svc := NewService(cosClient)
		require.NotNil(t, svc)

		// process the request
		actual, err := svc.GetAccount(ctx, rpcReq)
		require.Error(t, err)
		require.Nil(t, actual)
		cosClient.AssertExpectations(t)
	})
}
