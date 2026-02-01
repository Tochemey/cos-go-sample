package subscription

import (
	"context"

	"github.com/tochemey/cos-go-sample/app/storage"
)

type SubscriptionHandler struct {
	dataStore storage.Storage
}

func NewSubscriptionHandler(dataStore storage.Storage) *SubscriptionHandler {
	return &SubscriptionHandler{
		dataStore: dataStore,
	}
}

func (s SubscriptionHandler) HandleEvents(ctx context.Context, events []interface{}) error {
	panic("implement me")
}
