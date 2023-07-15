package writeside

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
	cospb "github.com/tochemey/cos-go-sample/gen/chief_of_state/v1"
	commands "github.com/tochemey/cos-go-sample/mocks/app/services/writeside/commands"
	events "github.com/tochemey/cos-go-sample/mocks/app/services/writeside/events"
	gopack "github.com/tochemey/gopack/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

type serviceSuite struct {
	suite.Suite
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestHandlerService(t *testing.T) {
	suite.Run(t, new(serviceSuite))
}

func (s *serviceSuite) TestNewHandlerService() {
	mockCommandHandler := new(commands.Dispatcher)
	mockEventHandler := new(events.Dispatcher)
	svc := NewHandlerService(mockCommandHandler, mockEventHandler)
	s.Assert().NotNil(svc)
}

func (s *serviceSuite) TestRegisterService() {
	// just for test coverage sake
	svc := &HandlerService{}
	// get in process grpc test server
	testServer, _ := gopack.TestServer(nil)
	s.Assert().NotPanics(func() { svc.RegisterService(testServer) })
}

func (s *serviceSuite) TestHandleCommand() {
	s.Run("With valid command", func() {
		// create the request context
		ctx := context.TODO()

		// create the command to handle
		actualCmd := new(pb.OpenAccount)
		// marshall the command as an any
		command, err := anypb.New(actualCmd)

		s.Assert().NoError(err)
		s.Assert().NotNil(command)

		// create the priorState
		actualState := new(pb.BankAccount)
		priorState, err := anypb.New(actualState)

		s.Assert().NoError(err)
		s.Assert().NotNil(priorState)

		// create the prior event meta
		priorEventMeta := new(cospb.MetaData)

		// create the handle command request
		request := &cospb.HandleCommandRequest{
			Command:        command,
			PriorState:     priorState,
			PriorEventMeta: priorEventMeta,
		}

		// create the resulting event when the command has been handled successfully
		event := new(pb.AccountOpened)

		// create the command handler and mock the HandleCommand method
		mockCommandHandler := new(commands.Dispatcher)
		mockCommandHandler.On("Dispatch", ctx, actualCmd, actualState, priorEventMeta).Return(event, nil)
		mockEventHandler := new(events.Dispatcher)

		// create an instance of the service with the mocks
		svc := NewHandlerService(mockCommandHandler, mockEventHandler)
		response, err := svc.HandleCommand(ctx, request)

		// assert the response and error
		s.Assert().NoError(err)
		s.Assert().NotNil(response)
		s.Assert().True(response.GetEvent().MessageIs(event))
		// make sure the command handler has been called
		mockCommandHandler.AssertExpectations(s.T())
	})
	s.Run("With invalid command", func() {
		// create the request context
		ctx := context.TODO()

		// mock the invalidCommand and event handler
		mockCommandHandler := new(commands.Dispatcher)
		mockEventHandler := new(events.Dispatcher)

		// create the priorState
		actualState := new(pb.BankAccount)
		priorState, err := anypb.New(actualState)

		s.Assert().NoError(err)
		s.Assert().NotNil(priorState)

		// create the prior event meta
		priorEventMeta := new(cospb.MetaData)
		// create the invalid invalidCommand
		invalidCommand := &anypb.Any{
			TypeUrl: "",
			Value:   nil,
		}

		// create the handle invalidCommand request
		request := &cospb.HandleCommandRequest{
			Command:        invalidCommand,
			PriorState:     priorState,
			PriorEventMeta: priorEventMeta,
		}

		// create an instance of the service with the mocks
		svc := NewHandlerService(mockCommandHandler, mockEventHandler)
		response, err := svc.HandleCommand(ctx, request)

		// assert the response and error
		s.Assert().Error(err)
		s.Assert().Nil(response)
	})
	s.Run("With invalid state", func() {
		// create the request context
		ctx := context.TODO()

		// mock the command and event handler
		mockCommandHandler := new(commands.Dispatcher)
		mockEventHandler := new(events.Dispatcher)

		// create the command to handle
		actualCmd := new(pb.OpenAccount)
		command, err := anypb.New(actualCmd)

		s.Assert().NoError(err)
		s.Assert().NotNil(command)

		// create the prior event meta
		priorEventMeta := new(cospb.MetaData)

		// create an invalid state
		invalidState := &anypb.Any{
			TypeUrl: "",
			Value:   nil,
		}

		// create the handle command request
		request := &cospb.HandleCommandRequest{
			Command:        command,
			PriorState:     invalidState,
			PriorEventMeta: priorEventMeta,
		}

		// create an instance of the service with the mocks
		svc := NewHandlerService(mockCommandHandler, mockEventHandler)
		response, err := svc.HandleCommand(ctx, request)

		// assert the response and error
		s.Assert().Error(err)
		s.Assert().Nil(response)
	})
	s.Run("With failed command handler", func() {
		// create the request context
		ctx := context.TODO()

		// create the command to handle
		actualCmd := new(pb.OpenAccount)
		// marshall the command as an anypb
		command, err := anypb.New(actualCmd)

		s.Assert().NoError(err)
		s.Assert().NotNil(command)

		// create the priorState
		actualState := new(pb.BankAccount)
		priorState, err := anypb.New(actualState)

		s.Assert().NoError(err)
		s.Assert().NotNil(priorState)

		// create the prior event meta
		priorEventMeta := new(cospb.MetaData)

		// create the handle command request
		request := &cospb.HandleCommandRequest{
			Command:        command,
			PriorState:     priorState,
			PriorEventMeta: priorEventMeta,
		}

		// create the command handler error
		handlerErr := status.Error(codes.InvalidArgument, "some-error")

		// create the command handler and mock the HandleCommand method
		mockCommandHandler := new(commands.Dispatcher)
		mockCommandHandler.On("Dispatch", ctx, actualCmd, actualState, priorEventMeta).Return(nil, handlerErr)
		mockEventHandler := new(events.Dispatcher)

		// create an instance of the service with the mocks
		svc := NewHandlerService(mockCommandHandler, mockEventHandler)
		response, err := svc.HandleCommand(ctx, request)

		// assert the response and error
		s.Assert().Error(err)
		s.Assert().Nil(response)

		// make sure the command handler has been called
		mockCommandHandler.AssertExpectations(s.T())
	})
	s.Run("With nil proto message event returned", func() {
		// create the request context
		ctx := context.TODO()

		// create the command to handle
		actualCmd := new(pb.OpenAccount)
		command, err := anypb.New(actualCmd)

		s.Assert().NoError(err)
		s.Assert().NotNil(command)

		// create the priorState
		actualState := new(pb.BankAccount)
		priorState, err := anypb.New(actualState)

		s.Assert().NoError(err)
		s.Assert().NotNil(priorState)

		// create the prior event meta
		priorEventMeta := new(cospb.MetaData)

		// create the handle command request
		request := &cospb.HandleCommandRequest{
			Command:        command,
			PriorState:     priorState,
			PriorEventMeta: priorEventMeta,
		}

		// create the resulting event when the command has been handled successfully
		var event *pb.AccountOpened

		// create the command handler and mock the HandleCommand method
		mockCommandHandler := new(commands.Dispatcher)
		mockCommandHandler.On("Dispatch", ctx, actualCmd, actualState, priorEventMeta).Return(event, nil)
		mockEventHandler := new(events.Dispatcher)

		// create an instance of the service with the mocks
		svc := NewHandlerService(mockCommandHandler, mockEventHandler)
		response, err := svc.HandleCommand(ctx, request)

		// assert the response and error
		s.Assert().NoError(err)
		s.Assert().NotNil(response)
		s.Assert().False(response.GetEvent().MessageIs(new(pb.AccountOpened)))
		// make sure the command handler has been called
		mockCommandHandler.AssertExpectations(s.T())
	})
	s.Run("With nil event returned", func() {
		// create the request context
		ctx := context.TODO()

		// create the command to handle
		actualCmd := new(pb.OpenAccount)
		command, err := anypb.New(actualCmd)

		s.Assert().NoError(err)
		s.Assert().NotNil(command)

		// create the priorState
		actualState := new(pb.BankAccount)
		priorState, err := anypb.New(actualState)

		s.Assert().NoError(err)
		s.Assert().NotNil(priorState)

		// create the prior event meta
		priorEventMeta := new(cospb.MetaData)

		// create the handle command request
		request := &cospb.HandleCommandRequest{
			Command:        command,
			PriorState:     priorState,
			PriorEventMeta: priorEventMeta,
		}

		// create the command handler and mock the HandleCommand method
		mockCommandHandler := new(commands.Dispatcher)
		mockCommandHandler.On("Dispatch", ctx, actualCmd, actualState, priorEventMeta).Return(nil, nil)
		mockEventHandler := new(events.Dispatcher)

		// create an instance of the service with the mocks
		svc := NewHandlerService(mockCommandHandler, mockEventHandler)
		response, err := svc.HandleCommand(ctx, request)

		// assert the response and error
		s.Assert().NoError(err)
		s.Assert().NotNil(response)
		s.Assert().False(response.GetEvent().MessageIs(new(pb.AccountOpened)))
		// make sure the command handler has been called
		mockCommandHandler.AssertExpectations(s.T())
	})
}

func (s *serviceSuite) TestHandleEvent() {
	s.Run("With valid event", func() {
		// create the request context
		ctx := context.TODO()

		// create the event to handle
		actualEvent := new(pb.AccountOpened)
		event, err := anypb.New(actualEvent)

		s.Assert().NoError(err)
		s.Assert().NotNil(event)

		// create the priorState
		actualState := new(pb.BankAccount)
		priorState, err := anypb.New(actualState)

		s.Assert().NoError(err)
		s.Assert().NotNil(priorState)

		// create resulting state
		resultingState := &pb.BankAccount{
			AccountId: "id-1",
		}

		// create the prior event meta
		priorEventMeta := new(cospb.MetaData)
		// create the command handler
		mockCommandHandler := new(commands.Dispatcher)
		// create the event handler and mock HandleEvent
		mockEventHandler := new(events.Dispatcher)
		mockEventHandler.On("Dispatch", ctx, actualEvent, actualState, priorEventMeta).Return(resultingState, nil)

		// create the handle command request
		request := &cospb.HandleEventRequest{
			Event:      event,
			PriorState: priorState,
			EventMeta:  priorEventMeta,
		}

		// create an instance of the service with the mocks
		svc := NewHandlerService(mockCommandHandler, mockEventHandler)
		response, err := svc.HandleEvent(ctx, request)

		// assert the response and error
		s.Assert().NoError(err)
		s.Assert().NotNil(response)
		s.Assert().True(response.GetResultingState().MessageIs(new(pb.BankAccount)))
		// make sure the event handler has been called
		mockEventHandler.AssertExpectations(s.T())
	})
	s.Run("With invalid event", func() {
		// create the request context
		ctx := context.TODO()

		// create an invalid event
		invalidEvent := &anypb.Any{
			TypeUrl: "",
			Value:   nil,
		}

		// create the priorState
		actualState := new(pb.BankAccount)
		priorState, err := anypb.New(actualState)

		s.Assert().NoError(err)
		s.Assert().NotNil(priorState)

		// create the prior event meta
		priorEventMeta := new(cospb.MetaData)
		// create the command handler
		mockCommandHandler := new(commands.Dispatcher)
		// create the event handler and mock HandleEvent
		mockEventHandler := new(events.Dispatcher)

		// create the handle command request
		request := &cospb.HandleEventRequest{
			Event:      invalidEvent,
			PriorState: priorState,
			EventMeta:  priorEventMeta,
		}

		// create an instance of the service with the mocks
		svc := NewHandlerService(mockCommandHandler, mockEventHandler)
		response, err := svc.HandleEvent(ctx, request)

		// assert the response and error
		s.Assert().Error(err)
		s.Assert().Nil(response)
	})
	s.Run("With invalid state", func() {
		// create the request context
		ctx := context.TODO()

		// create the event to handle
		actualEvent := new(pb.AccountOpened)
		event, err := anypb.New(actualEvent)

		s.Assert().NoError(err)
		s.Assert().NotNil(event)

		// create an invalid state
		invalidState := &anypb.Any{
			TypeUrl: "",
			Value:   nil,
		}

		// create the prior event meta
		priorEventMeta := new(cospb.MetaData)
		// create the command handler
		mockCommandHandler := new(commands.Dispatcher)
		// create the event handler and mock HandleEvent
		mockEventHandler := new(events.Dispatcher)

		// create the handle command request
		request := &cospb.HandleEventRequest{
			Event:      event,
			PriorState: invalidState,
			EventMeta:  priorEventMeta,
		}

		// create an instance of the service with the mocks
		svc := NewHandlerService(mockCommandHandler, mockEventHandler)
		response, err := svc.HandleEvent(ctx, request)

		// assert the response and error
		s.Assert().Error(err)
		s.Assert().Nil(response)
	})
	s.Run("With failed event handler", func() {
		// create the request context
		ctx := context.TODO()

		// create the event to handle
		actualEvent := new(pb.AccountOpened)
		event, err := anypb.New(actualEvent)
		s.Assert().NoError(err)
		s.Assert().NotNil(event)

		// create the priorState
		actualState := new(pb.BankAccount)
		priorState, err := anypb.New(actualState)

		s.Assert().NoError(err)
		s.Assert().NotNil(priorState)

		// create the prior event meta
		priorEventMeta := new(cospb.MetaData)
		// create the command handler
		mockCommandHandler := new(commands.Dispatcher)
		// create the event handler and mock HandleEvent
		mockEventHandler := new(events.Dispatcher)
		mockEventHandler.
			On("Dispatch", ctx, actualEvent, actualState, priorEventMeta).
			Return(nil, status.Error(codes.Internal, "some-error"))

		// create the handle command request
		request := &cospb.HandleEventRequest{
			Event:      event,
			PriorState: priorState,
			EventMeta:  priorEventMeta,
		}

		// create an instance of the service with the mocks
		svc := NewHandlerService(mockCommandHandler, mockEventHandler)
		response, err := svc.HandleEvent(ctx, request)

		// assert the response and error
		s.Assert().Error(err)
		s.Assert().Nil(response)

		// make sure the event handler has been called
		mockEventHandler.AssertExpectations(s.T())
	})
	s.Run("with nil(proto message) resulting state returned", func() {
		// create the request context
		ctx := context.TODO()

		// create the event to handle
		actualEvent := new(pb.AccountOpened)
		event, err := anypb.New(actualEvent)

		s.Assert().NoError(err)
		s.Assert().NotNil(event)

		// create the priorState
		actualState := new(pb.BankAccount)
		priorState, err := anypb.New(actualState)

		s.Assert().NoError(err)
		s.Assert().NotNil(priorState)

		var newState *pb.BankAccount
		// create the prior event meta
		priorEventMeta := new(cospb.MetaData)
		// create the command handler
		mockCommandHandler := new(commands.Dispatcher)
		// create the event handler and mock HandleEvent
		mockEventHandler := new(events.Dispatcher)
		mockEventHandler.On("Dispatch", ctx, actualEvent, actualState, priorEventMeta).Return(newState, nil)

		// create the handle command request
		request := &cospb.HandleEventRequest{
			Event:      event,
			PriorState: priorState,
			EventMeta:  priorEventMeta,
		}

		// create an instance of the service with the mocks
		svc := NewHandlerService(mockCommandHandler, mockEventHandler)
		response, err := svc.HandleEvent(ctx, request)

		// assert the response and error
		s.Assert().NoError(err)
		s.Assert().NotNil(response)
		// make sure the event handler has been called
		mockEventHandler.AssertExpectations(s.T())
	})
}
