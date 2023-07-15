package events

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
	"google.golang.org/protobuf/proto"
)

func TestAccountOpened(t *testing.T) {
	ctx := context.TODO()
	accountID := "account-1"
	accountBal := 50.00
	accountOwner := "John Doe"
	amount := 50.00

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

	actual, err := accountOpened(ctx, event)
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.IsType(t, new(pb.BankAccount), actual)
	assert.True(t, proto.Equal(expected, actual))
}
