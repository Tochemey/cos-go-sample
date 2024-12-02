package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
)

func TestOpenAccount(t *testing.T) {
	ctx := context.TODO()
	accountID := "account-1"
	amount := 50.00

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

	// perform the credit account command handling
	actual, err := openAccount(ctx, command)
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.IsType(t, new(pb.AccountOpened), actual)
	assert.True(t, proto.Equal(expected, actual))
}
