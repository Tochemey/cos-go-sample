package events

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
	cospb "github.com/tochemey/cos-go-sample/gen/chief_of_state/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestNewDispatcher(t *testing.T) {
	dispatcher := NewDispatcher()
	assert.NotNil(t, dispatcher)
	var p interface{} = dispatcher
	_, ok := p.(Dispatcher)
	assert.True(t, ok)
}

func TestDispatch(t *testing.T) {
	t.Run("with nil event", func(t *testing.T) {
		// define a context
		ctx := context.TODO()
		// define the prior state
		priorState := new(pb.BankAccount)
		// define the CoS meta
		cosMeta := &cospb.MetaData{}
		actual, err := NewDispatcher().Dispatch(ctx, nil, priorState, cosMeta)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.EqualError(t, err, errEventNotDefined.Error())
	})
	t.Run("with unknown event", func(t *testing.T) {
		// define a context
		ctx := context.TODO()
		// define the prior state
		priorState := new(pb.BankAccount)
		// define the CoS meta
		cosMeta := &cospb.MetaData{}
		event := &emptypb.Empty{}
		actual, err := NewDispatcher().Dispatch(ctx, event, priorState, cosMeta)
		assert.Error(t, err)
		assert.EqualError(t, err, errUnhandledEvent(event).Error())
		assert.Nil(t, actual)
	})
	t.Run("With AccountOpened event", func(t *testing.T) {
		ctx := context.TODO()
		accountID := "account-1"
		accountBal := 50.00
		accountOwner := "John Doe"
		amount := 50.00

		// create the prior state
		priorState := &pb.BankAccount{}

		// create the event
		event := &pb.AccountOpened{
			AccountId:    accountID,
			Balance:      amount,
			AccountOwner: accountOwner,
		}

		expected := &pb.BankAccount{
			AccountId:      accountID,
			AccountBalance: accountBal,
			AccountOwner:   accountOwner,
			IsClosed:       false,
		}

		// create the cos prior meta
		priorMeta := &cospb.MetaData{EntityId: accountID}

		// create the instance of the dispatcher
		dispatcher := NewDispatcher()

		// perform the credit account command handling
		actual, err := dispatcher.Dispatch(ctx, event, priorState, priorMeta)
		require.NoError(t, err)
		require.NotNil(t, actual)
		require.IsType(t, new(pb.BankAccount), actual)
		assert.True(t, proto.Equal(expected, actual))
	})
	t.Run("With AccountCredited event", func(t *testing.T) {
		ctx := context.TODO()
		accountID := "account-1"
		accountBal := 150.55
		accountOwner := "John Doe"
		amount := 50.00

		// create the prior state
		priorState := &pb.BankAccount{
			AccountId:      accountID,
			AccountBalance: accountBal,
			AccountOwner:   accountOwner,
			IsClosed:       false,
		}

		// create the event
		event := &pb.AccountCredited{
			AccountId: accountID,
			Amount:    amount,
		}

		expected := &pb.BankAccount{
			AccountId:      accountID,
			AccountBalance: accountBal + amount,
			AccountOwner:   accountOwner,
			IsClosed:       false,
		}

		// create the cos prior meta
		priorMeta := &cospb.MetaData{EntityId: accountID}

		// create the instance of the dispatcher
		dispatcher := NewDispatcher()

		// perform the credit account command handling
		actual, err := dispatcher.Dispatch(ctx, event, priorState, priorMeta)
		require.NoError(t, err)
		require.NotNil(t, actual)
		require.IsType(t, new(pb.BankAccount), actual)
		assert.True(t, proto.Equal(expected, actual))
	})
	t.Run("with AccountDebited event", func(t *testing.T) {
		ctx := context.TODO()
		accountID := "account-1"
		accountBal := 150.55
		accountOwner := "John Doe"
		amount := 50.00

		// create the prior state
		priorState := &pb.BankAccount{
			AccountId:      accountID,
			AccountBalance: accountBal,
			AccountOwner:   accountOwner,
			IsClosed:       false,
		}

		// create the event
		event := &pb.AccountDebited{
			AccountId: accountID,
			Amount:    amount,
		}

		expected := &pb.BankAccount{
			AccountId:      accountID,
			AccountBalance: accountBal - amount,
			AccountOwner:   accountOwner,
			IsClosed:       false,
		}

		// create the cos prior meta
		priorMeta := &cospb.MetaData{EntityId: accountID}

		// create the instance of the dispatcher
		dispatcher := NewDispatcher()

		// perform the credit account command handling
		actual, err := dispatcher.Dispatch(ctx, event, priorState, priorMeta)
		require.NoError(t, err)
		require.NotNil(t, actual)
		require.IsType(t, new(pb.BankAccount), actual)
		assert.True(t, proto.Equal(expected, actual))
	})
}
