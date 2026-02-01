package subscription

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	gopack "github.com/tochemey/gopack/grpc"
	"google.golang.org/protobuf/proto"

	cospb "github.com/tochemey/cos-go-sample/gen/chief_of_state/v1"
)

// Manager manages Chief of State event subscriptions.
// It subscribes to all events on Start and unsubscribes on Stop.
type Manager struct {
	cosHost   string
	cosPort   int
	handler   *Handler
	conn      interface{ Close() error }
	cosClient cospb.ChiefOfStateServiceClient
	subID     string
	subMux    sync.RWMutex
	stopCh    chan struct{}
	stopOnce  sync.Once
}

// NewManager creates a new subscription manager.
func NewManager(cosHost string, cosPort int, handler *Handler) (*Manager, error) {
	conn, err := gopack.DefaultConn(fmt.Sprintf("%s:%d", cosHost, cosPort))
	if err != nil {
		return nil, err
	}
	return &Manager{
		cosHost:   cosHost,
		cosPort:   cosPort,
		handler:   handler,
		conn:      conn,
		cosClient: cospb.NewChiefOfStateServiceClient(conn),
		stopCh:    make(chan struct{}),
	}, nil
}

// Start begins subscribing to all CoS events and forwards them to the handler.
// It starts a background goroutine and returns immediately.
// Call Stop to unsubscribe and clean up.
func (m *Manager) Start(ctx context.Context) error {
	subID := uuid.New().String()
	req := &cospb.SubscribeAllRequest{SubscriptionId: subID}

	stream, err := m.cosClient.SubscribeAll(ctx, req)
	if err != nil {
		return fmt.Errorf("subscribe all: %w", err)
	}

	m.subMux.Lock()
	m.subID = subID
	m.subMux.Unlock()

	go m.receiveLoop(ctx, stream)
	return nil
}

func (m *Manager) receiveLoop(ctx context.Context, stream interface {
	Recv() (*cospb.SubscribeAllResponse, error)
}) {
	for {
		select {
		case <-m.stopCh:
			return
		case <-ctx.Done():
			return
		default:
			resp, err := stream.Recv()
			if err != nil {
				select {
				case <-m.stopCh:
				default:
					// Stream closed (e.g. context cancelled during shutdown)
				}
				return
			}
			if resp == nil {
				continue
			}
			m.subMux.Lock()
			if resp.SubscriptionId != "" {
				m.subID = resp.SubscriptionId
			}
			m.subMux.Unlock()
			events := m.convertToEvents(resp)
			if len(events) > 0 && m.handler != nil {
				_ = m.handler.HandleEvents(ctx, events)
			}
		}
	}
}

// convertToEvents converts a SubscribeAllResponse to events for the handler.
func (m *Manager) convertToEvents(resp *cospb.SubscribeAllResponse) []any {
	if resp.GetEvent() == nil {
		return nil
	}

	event, err := resp.Event.UnmarshalNew()
	if err != nil {
		return []any{&UnknownEvent{TypeURL: resp.Event.GetTypeUrl(), Raw: resp.Event}}
	}

	item := &Event{
		Event:          event,
		ResultingState: nil,
		Meta:           resp.GetMeta(),
	}

	if resp.GetResultingState() != nil {
		state, err := resp.ResultingState.UnmarshalNew()
		if err == nil {
			item.ResultingState = state
		}
	}

	return []any{item}
}

// Stop unsubscribes from all events and closes the connection.
func (m *Manager) Stop(ctx context.Context) error {
	var err error
	m.stopOnce.Do(func() {
		close(m.stopCh)
		m.subMux.RLock()
		subID := m.subID
		m.subMux.RUnlock()
		if subID != "" {
			unsubCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, unsubErr := m.cosClient.UnsubscribeAll(unsubCtx, &cospb.UnsubscribeAllRequest{SubscriptionId: subID})
			if unsubErr != nil {
				err = fmt.Errorf("unsubscribe all: %w", unsubErr)
			}
		}
		if m.conn != nil {
			_ = m.conn.Close()
		}
	})
	return err
}

// Event represents a single event from the subscription stream.
type Event struct {
	Event          proto.Message
	ResultingState proto.Message
	Meta           *cospb.MetaData
}

// UnknownEvent represents an event that could not be unmarshaled.
type UnknownEvent struct {
	TypeURL string
	Raw     proto.Message
}
