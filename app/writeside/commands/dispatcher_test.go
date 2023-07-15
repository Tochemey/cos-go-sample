package commands

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
	t.Run("with nil command", func(t *testing.T) {
		// create a context
		ctx := context.TODO()
		// create the prior state
		priorState := &pb.BankAccount{}
		// create the CoS meta
		cosMeta := &cospb.MetaData{}
		handler := &dispatcher{}
		resultingState, err := handler.Dispatch(ctx, nil, priorState, cosMeta)
		assert.Error(t, err)
		assert.Nil(t, resultingState)
		assert.EqualError(t, err, errCommandNotDefined.Error())
	})
	t.Run("with unknown command", func(t *testing.T) {
		// create a context
		ctx := context.TODO()
		// create the prior state
		priorState := &pb.BankAccount{}
		// create the CoS meta
		cosMeta := &cospb.MetaData{}
		dispatch := &dispatcher{}
		command := &emptypb.Empty{}
		resultingState, err := dispatch.Dispatch(ctx, command, priorState, cosMeta)
		assert.Nil(t, resultingState)
		assert.Error(t, err)
		assert.EqualError(t, err, errUnhandledCommand(command).Error())
	})
	t.Run("With CreditAccount command", func(t *testing.T) {
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

		// create the command
		command := &pb.CreditAccount{
			AccountId: accountID,
			Amount:    amount,
		}

		// create the expected outcome
		expected := &pb.AccountCredited{
			AccountId: accountID,
			Amount:    amount,
		}

		// create the cos prior meta
		priorMeta := &cospb.MetaData{EntityId: accountID}

		// create the instance of the dispatcher
		dispatcher := NewDispatcher()

		// perform the credit account command handling
		actual, err := dispatcher.Dispatch(ctx, command, priorState, priorMeta)
		require.NoError(t, err)
		require.NotNil(t, actual)
		require.IsType(t, new(pb.AccountCredited), actual)
		assert.True(t, proto.Equal(expected, actual))
	})
	t.Run("With OpenAccount command", func(t *testing.T) {
		ctx := context.TODO()
		accountID := "account-1"
		amount := 50.00

		// create the prior state
		priorState := &pb.BankAccount{}

		// create the command
		command := &pb.OpenAccount{
			AccountId:      accountID,
			OpeningBalance: amount,
		}

		// create the expected outcome
		expected := &pb.AccountOpened{
			AccountId: accountID,
			Balance:   amount,
		}

		// create the cos prior meta
		priorMeta := &cospb.MetaData{EntityId: accountID, RevisionNumber: 2}

		// create the instance of the dispatcher
		dispatcher := NewDispatcher()

		// perform the credit account command handling
		actual, err := dispatcher.Dispatch(ctx, command, priorState, priorMeta)
		require.NoError(t, err)
		require.NotNil(t, actual)
		require.IsType(t, new(pb.AccountOpened), actual)
		assert.True(t, proto.Equal(expected, actual))
	})
	t.Run("With DebitAccount command", func(t *testing.T) {
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

		// create the command
		command := &pb.DebitAccount{
			AccountId: accountID,
			Amount:    amount,
		}

		// create the expected outcome
		expected := &pb.AccountDebited{
			AccountId: accountID,
			Amount:    amount,
		}
		// create the cos prior meta
		priorMeta := &cospb.MetaData{EntityId: accountID, RevisionNumber: 2}

		// create the instance of the dispatcher
		dispatcher := NewDispatcher()

		// perform the credit account command handling
		actual, err := dispatcher.Dispatch(ctx, command, priorState, priorMeta)
		require.NoError(t, err)
		require.NotNil(t, actual)
		require.IsType(t, new(pb.AccountDebited), actual)
		assert.True(t, proto.Equal(expected, actual))
	})
}
