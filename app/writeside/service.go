package writeside

import (
	"context"

	"github.com/tochemey/cos-go-sample/app/writeside/commands"
	"github.com/tochemey/cos-go-sample/app/writeside/events"

	"github.com/pkg/errors"
	"github.com/tochemey/cos-go-sample/app/cos"
	cospb "github.com/tochemey/cos-go-sample/gen/chief_of_state/v1"
	"github.com/tochemey/gopack/log/zapl"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
)

// HandlerService is an implementation of the CoS WriteSide handler interface
type HandlerService struct {
	commandsDispatcher commands.Dispatcher
	eventsDispatcher   events.Dispatcher
}

// enforce compilation error when the HandlerService does not fully implement the WriteSideHandlerServiceServer
// interface
var _ cospb.WriteSideHandlerServiceServer = (*HandlerService)(nil)

// NewHandlerService creates a new instance of HandlerService
func NewHandlerService(commandsDispatcher commands.Dispatcher, eventsDispatcher events.Dispatcher) *HandlerService {
	// create the service object and set the commands and events handler
	return &HandlerService{
		commandsDispatcher: commandsDispatcher,
		eventsDispatcher:   eventsDispatcher,
	}
}

// HandleCommand accepts commands and returns events for CoS write handler
func (s HandlerService) HandleCommand(ctx context.Context, request *cospb.HandleCommandRequest) (*cospb.HandleCommandResponse, error) {
	// set the log with the context
	log := zapl.WithContext(ctx)

	// unpacking the command
	cmd, err := request.GetCommand().UnmarshalNew()
	if err != nil {
		err = errors.Wrapf(err, "failed to unpack command:(%s)", request.GetCommand().GetTypeUrl())
		log.Error(err)
		return nil, err
	}

	// unpacking the state
	priorState, err := cos.UnmarshalState(request.GetPriorState())
	if err != nil {
		err = errors.Wrapf(err, "failed to unpack state:(%s)", request.GetPriorState().GetTypeUrl())
		log.Error(err)
		return nil, err
	}

	event, err := s.commandsDispatcher.Dispatch(ctx, cmd, priorState, request.GetPriorEventMeta())
	if err != nil {
		err = errors.Wrapf(err, "failed to handle command:(%s)", cmd.ProtoReflect().Descriptor().FullName())
		log.Error(err)
		return nil, err
	}

	// prepare response object
	response := &cospb.HandleCommandResponse{}

	// if there is an event, inject as Any, else leave nil to signal a no-op to COS
	if event != nil {
		eventAny, err := anypb.New(event)
		if err != nil {
			err = errors.Wrapf(err, "failed to pack event:(%s) as any proto message",
				event.ProtoReflect().Descriptor().FullName())
			log.Error(err)
			return nil, err
		}

		// this means the original event was nil, aka no-op
		// prior nil check for response might not return nil in a typed nil vs. nil situation, so we check here as well.
		if eventAny.GetValue() == nil {
			return response, nil
		}

		// set the event and return
		response.Event = eventAny
	}

	return response, nil
}

// HandleEvent accepts events and returns new states for CoS write handler
func (s HandlerService) HandleEvent(ctx context.Context, request *cospb.HandleEventRequest) (*cospb.HandleEventResponse, error) {
	// set the log with the context
	log := zapl.WithContext(ctx)
	event, err := request.GetEvent().UnmarshalNew()
	if err != nil {
		err = errors.Wrapf(err, "failed to unpack event:(%s)", request.GetEvent().GetTypeUrl())
		log.Error(err)
		return nil, err
	}

	// unpack the prior state
	state, err := cos.UnmarshalState(request.GetPriorState())
	// handle the error
	if err != nil {
		err = errors.Wrapf(err, "failed to unpack state:(%s)", request.GetPriorState().GetTypeUrl())
		log.Error(err)
		return nil, err
	}

	// handle the event
	resultingState, err := s.eventsDispatcher.Dispatch(ctx, event, state, request.GetEventMeta())
	// handle the error
	if err != nil {
		err = errors.Wrapf(err, "failed to handle event:(%s)", event.ProtoReflect().Descriptor().FullName())
		log.Error(err)
		return nil, err
	}

	// pack the resulting state as any proto message
	resultingStateAny, err := anypb.New(resultingState)
	if err != nil {
		err = errors.Wrapf(err, "failed to pack resulting state:(%s) as any proto message",
			resultingState.ProtoReflect().Descriptor().FullName())
		log.Error(err)
		return nil, err
	}

	// in case we have empty resulting state
	if resultingStateAny.GetValue() == nil || len(resultingStateAny.GetValue()) == 0 {
		return new(cospb.HandleEventResponse), nil
	}

	return &cospb.HandleEventResponse{ResultingState: resultingStateAny}, nil
}

// RegisterService registers the gRPC api
func (s HandlerService) RegisterService(sv *grpc.Server) {
	cospb.RegisterWriteSideHandlerServiceServer(sv, s)
}
