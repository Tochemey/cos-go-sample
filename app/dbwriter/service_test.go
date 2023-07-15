package dbwriter

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
	cospb "github.com/tochemey/cos-go-sample/gen/chief_of_state/v1"
	mocks "github.com/tochemey/cos-go-sample/mocks/app/storage"
	gopack "github.com/tochemey/gopack/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestNewService(t *testing.T) {
	t.Run("With happy path", func(t *testing.T) {
		dataStore := new(mocks.Storage)
		svc, err := NewService(dataStore)
		assert.NoError(t, err)
		assert.NotNil(t, svc)

		// assert the type of svc
		assert.IsType(t, &Service{}, svc)

		// assert that svc implement cospb.ReadSideHandlerServiceServer
		var p interface{} = svc
		_, ok := p.(cospb.ReadSideHandlerServiceServer)
		assert.True(t, ok)
	})
	t.Run("With data store not set", func(t *testing.T) {
		svc, err := NewService(nil)
		assert.Error(t, err)
		assert.EqualError(t, err, "the dataStore is not defined")
		assert.Nil(t, svc)
	})
}

func TestRegisterService(t *testing.T) {
	fn := func() {
		svc := &Service{}
		// get in process grpc test server
		testServer, _ := gopack.TestServer(nil)
		svc.RegisterService(testServer)
	}

	assert.NotPanics(t, fn)
}

func TestHandleReadSide(t *testing.T) {
	t.Run("With happy path", func(t *testing.T) {
		ctx := context.TODO()
		accountID := "account-1"
		state := &pb.BankAccount{AccountId: accountID}
		// pack the state into any pb
		anyState, err := anypb.New(state)
		require.NoError(t, err)
		require.NotNil(t, anyState)
		// create mocks
		dataStore := new(mocks.Storage)
		dataStore.On("PersistAccount", ctx, mock.MatchedBy(func(in *pb.BankAccount) bool {
			return proto.Equal(in, state)
		})).Return(nil)

		svc, err := NewService(dataStore)
		require.NoError(t, err)
		require.NotNil(t, svc)

		// create the read side request with the relevant needed info
		req := &cospb.HandleReadSideRequest{State: anyState}
		// handle the read side request
		resp, err := svc.HandleReadSide(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.GetSuccessful())
		dataStore.AssertExpectations(t)
	})
	t.Run("With nil state", func(t *testing.T) {
		ctx := context.TODO()
		// create mocks
		dataStore := new(mocks.Storage)
		svc, err := NewService(dataStore)
		require.NoError(t, err)
		require.NotNil(t, svc)

		// create the read side request with the relevant needed info
		req := &cospb.HandleReadSideRequest{}
		// handle the read side request
		resp, err := svc.HandleReadSide(ctx, req)
		assert.Error(t, err)
		assert.EqualError(t, err, "rpc error: code = Internal desc = the account state is not set")
		assert.Nil(t, resp)
		assert.False(t, resp.GetSuccessful())
		dataStore.AssertNotCalled(t, "PersistAccount")
	})
	t.Run("With wrong state", func(t *testing.T) {
		ctx := context.TODO()
		// pack the state into any pb
		anyState, err := anypb.New(wrapperspb.String("not a valid state"))
		require.NoError(t, err)
		require.NotNil(t, anyState)
		// create mocks
		dataStore := new(mocks.Storage)

		svc, err := NewService(dataStore)
		require.NoError(t, err)
		require.NotNil(t, svc)

		// create the read side request with the relevant needed info
		req := &cospb.HandleReadSideRequest{State: anyState}
		// handle the read side request
		resp, err := svc.HandleReadSide(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.False(t, resp.GetSuccessful())
		dataStore.AssertNotCalled(t, "PersistAccount")
	})
	t.Run("with dataStore failure", func(t *testing.T) {
		ctx := context.TODO()
		accountID := "account-1"
		state := &pb.BankAccount{AccountId: accountID}
		// pack the state into any pb
		anyState, err := anypb.New(state)
		require.NoError(t, err)
		require.NotNil(t, anyState)
		// create mocks
		dataStore := new(mocks.Storage)
		dataStore.On("PersistAccount", ctx, mock.MatchedBy(func(in *pb.BankAccount) bool {
			return proto.Equal(in, state)
		})).Return(errors.New("failed"))

		svc, err := NewService(dataStore)
		require.NoError(t, err)
		require.NotNil(t, svc)

		// create the read side request with the relevant needed info
		req := &cospb.HandleReadSideRequest{State: anyState}
		// handle the read side request
		resp, err := svc.HandleReadSide(ctx, req)
		assert.Error(t, err)
		assert.EqualError(t, err, "rpc error: code = Internal desc = failed to persist account into the data store: failed")
		assert.Nil(t, resp)
		assert.False(t, resp.GetSuccessful())
		dataStore.AssertExpectations(t)
	})
}
