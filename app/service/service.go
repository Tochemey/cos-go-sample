package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/tochemey/cos-go-sample/app/cos"
	pb "github.com/tochemey/cos-go-sample/gen/accounts/v1"
	"github.com/tochemey/gopack/log/zapl"
	"google.golang.org/grpc"
)

// Service implements the application service interface
type Service struct {
	cosClient cos.Client
}

// enforce compilation error when Service does not implement fully the
// BankAccountServiceServer interface
var _ pb.BankAccountServiceServer = &Service{}

// NewService creates an instance of api
func NewService(cosClient cos.Client) *Service {
	return &Service{
		cosClient,
	}
}

// OpenAccount helps open a bank account. When the request is successful the newly created account object is returned in the response.
// In case of error a gRPC error is returned. For more information refer to https://www.grpc.io/docs/guides/error/
func (s *Service) OpenAccount(ctx context.Context, request *pb.OpenAccountRequest) (*pb.OpenAccountResponse, error) {
	// get context log
	log := zapl.WithContext(ctx)

	// let us generate the account id or use it
	accountID := request.GetAccountId()
	if accountID == "" {
		accountID = uuid.NewString()
	}

	// let us create the command to send to CoS
	command := &pb.OpenAccount{
		AccountId:      accountID,
		AccountOwner:   request.GetAccountOwner(),
		OpeningBalance: request.GetBalance(),
	}

	// send the command to CoS
	state, _, err := s.cosClient.ProcessCommand(ctx, accountID, command)
	// handle the error
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &pb.OpenAccountResponse{Account: state}, nil
}

// DebitAccount sends a debit account request to the service. When the request is successful the debited account with the new balance is returned in the response.
// In case of error a gRPC error is returned. For more information refer to https://www.grpc.io/docs/guides/error/
func (s *Service) DebitAccount(ctx context.Context, request *pb.DebitAccountRequest) (*pb.DebitAccountResponse, error) {
	// get context log
	log := zapl.WithContext(ctx)

	// create the debit command
	command := &pb.DebitAccount{
		AccountId: request.GetAccountId(),
		Amount:    request.GetAmount(),
	}

	// send the request to CoS
	state, _, err := s.cosClient.ProcessCommand(ctx, request.GetAccountId(), command)
	// handle the error
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &pb.DebitAccountResponse{Account: state}, nil
}

// CreditAccount sends a credit account request to the service. When the request is successful the newly credited account with the new balance is returned in the response.
// In case of error a gRPC error is returned. For more information refer to https://www.grpc.io/docs/guides/error/
func (s *Service) CreditAccount(ctx context.Context, request *pb.CreditAccountRequest) (*pb.CreditAccountResponse, error) {
	// get context log
	log := zapl.WithContext(ctx)

	// create the command to send to CoS
	command := &pb.CreditAccount{
		AccountId: request.GetAccountId(),
		Amount:    request.GetAmount(),
	}

	// send the command to CoS
	state, _, err := s.cosClient.ProcessCommand(ctx, request.GetAccountId(), command)
	// handle the error
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &pb.CreditAccountResponse{Account: state}, nil
}

// GetAccount returns a given account information. When the request is successful the account info is returned in the response.
// In case of error a gRPC error is returned. For more information refer to https://www.grpc.io/docs/guides/error/
func (s *Service) GetAccount(ctx context.Context, request *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	// get context log
	log := zapl.WithContext(ctx)
	// let us get the current state from CoS. At this stage it makes sense to fetch the current state
	state, _, err := s.cosClient.GetState(ctx, request.GetAccountId())
	// handle the error
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &pb.GetAccountResponse{Account: state}, nil
}

// RegisterService registers the gRPC api
func (s *Service) RegisterService(sv *grpc.Server) {
	pb.RegisterBankAccountServiceServer(sv, s)
}
