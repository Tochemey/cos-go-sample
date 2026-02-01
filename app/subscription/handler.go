package subscription

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/tochemey/cos-go-sample/app/log"
	"github.com/tochemey/cos-go-sample/app/storage"
	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
)

type Handler struct {
	dataStore storage.Storage
}

func NewSubscriptionHandler(dataStore storage.Storage) *Handler {
	return &Handler{
		dataStore: dataStore,
	}
}

func (s *Handler) HandleEvents(ctx context.Context, events []any) error {
	logger := log.WithContext(ctx)
	for _, e := range events {
		evt, ok := e.(*Event)
		if !ok {
			if ue, ok := e.(*UnknownEvent); ok {
				logger.Infof("received unknown event: type=%s", ue.TypeURL)
			}
			continue
		}

		entityID := ""
		if evt.Meta != nil {
			entityID = evt.Meta.EntityId
		}

		logger.Infof("event received: entity_id=%s revision=%d type=%s",
			entityID, evt.Meta.GetRevisionNumber(), proto.MessageName(evt.Event))

		switch v := evt.Event.(type) {
		case *pb.AccountOpened:
			logger.Infof("  AccountOpened: account_id=%s owner=%s balance=%.2f",
				v.AccountId, v.AccountOwner, v.Balance)
		case *pb.AccountCredited:
			logger.Infof("  AccountCredited: account_id=%s amount=%.2f",
				v.AccountId, v.Amount)
		case *pb.AccountDebited:
			logger.Infof("  AccountDebited: account_id=%s amount=%.2f",
				v.AccountId, v.Amount)
		default:
			logger.Infof("  event: %+v", evt.Event)
		}
	}
	return nil
}
