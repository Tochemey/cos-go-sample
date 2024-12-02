package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
)

func TestDebitAccount(t *testing.T) {
	t.Run("With happy path", func(t *testing.T) {
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

		// perform the credit account command handling
		actual, err := debitAccount(ctx, command, priorState)
		require.NoError(t, err)
		require.NotNil(t, actual)
		require.IsType(t, new(pb.AccountDebited), actual)
		assert.True(t, proto.Equal(expected, actual))
	})
	t.Run("With prior state not defined", func(t *testing.T) {
		ctx := context.TODO()
		accountID := "account-1"
		amount := 50.00

		// create the prior state
		priorState := &pb.BankAccount{}

		// create the command
		command := &pb.DebitAccount{
			AccountId: accountID,
			Amount:    amount,
		}

		// perform the credit account command handling
		actual, err := debitAccount(ctx, command, priorState)
		require.Error(t, err)
		assert.EqualError(t, err, errMissingPriorState.Error())
		require.Nil(t, actual)
	})
	t.Run("With mismatch account id in command and prior state", func(t *testing.T) {
		ctx := context.TODO()
		accountID := "account-1"
		mismatchAccountID := "mismatch-1"
		amount := 50.00

		// create the prior state
		priorState := &pb.BankAccount{
			AccountId: mismatchAccountID,
		}

		// create the command
		command := &pb.DebitAccount{
			AccountId: accountID,
			Amount:    amount,
		}

		// perform the credit account command handling
		actual, err := debitAccount(ctx, command, priorState)
		require.Error(t, err)
		assert.EqualError(t, err, errCommandSentToWrongEntity.Error())
		require.Nil(t, actual)
	})
	t.Run("With insufficient balance", func(t *testing.T) {
		ctx := context.TODO()
		accountID := "account-1"
		accountBal := 20.00
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

		expectedErr := status.Error(codes.InvalidArgument, "insufficient balance")
		// perform the credit account command handling
		actual, err := debitAccount(ctx, command, priorState)
		require.Error(t, err)
		require.Nil(t, actual)
		assert.EqualError(t, err, expectedErr.Error())
	})
}
