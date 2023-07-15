package events

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
	"google.golang.org/protobuf/proto"
)

func TestAccountCredited(t *testing.T) {
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

	actual, err := accountCredited(ctx, event, priorState)
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.IsType(t, new(pb.BankAccount), actual)
	assert.True(t, proto.Equal(expected, actual))
}
