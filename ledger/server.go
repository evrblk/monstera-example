package ledger

import (
	"context"
	"log"
	"math/rand/v2"
	"time"

	"github.com/evrblk/monstera-example/ledger/corepb"
	"github.com/evrblk/monstera-example/ledger/gatewaypb"
	monsterax "github.com/evrblk/monstera/x"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LedgerServiceApiServer struct {
	gatewaypb.UnimplementedLedgerServiceApiServer

	coreApiClient LedgerServiceCoreApi
}

func (s *LedgerServiceApiServer) Close() {
	log.Println("Stopping ApiServer...")
}

func (s *LedgerServiceApiServer) CreateAccount(ctx context.Context, request *gatewaypb.CreateAccountRequest) (*gatewaypb.CreateAccountResponse, error) {
	now := time.Now()

	res1, err := s.coreApiClient.CreateAccount(ctx, &corepb.CreateAccountRequest{
		AccountId: rand.Uint64(),
		Now:       now.UnixNano(),
	})
	if err != nil {
		return nil, monsterax.ErrorToGRPC(err)
	}

	return &gatewaypb.CreateAccountResponse{
		Account: accountToFront(res1.Account),
	}, nil
}

func (s *LedgerServiceApiServer) GetAccount(ctx context.Context, request *gatewaypb.GetAccountRequest) (*gatewaypb.GetAccountResponse, error) {
	accountId, err := DecodeAccountId(request.AccountId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid account id")
	}

	resp1, err := s.coreApiClient.GetAccount(ctx, &corepb.GetAccountRequest{
		AccountId: accountId,
	})
	if err != nil {
		return nil, monsterax.ErrorToGRPC(err)
	}

	return &gatewaypb.GetAccountResponse{
		Account: accountToFront(resp1.Account),
	}, nil
}

func (s *LedgerServiceApiServer) GetTransaction(ctx context.Context, request *gatewaypb.GetTransactionRequest) (*gatewaypb.GetTransactionResponse, error) {
	transactionId, err := DecodeTransactionId(request.TransactionId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid transaction id")
	}

	resp1, err := s.coreApiClient.GetTransaction(ctx, &corepb.GetTransactionRequest{
		TransactionId: transactionId,
	})
	if err != nil {
		return nil, monsterax.ErrorToGRPC(err)
	}

	return &gatewaypb.GetTransactionResponse{
		Transaction: transactionToFront(resp1.Transaction),
	}, nil
}

func (s *LedgerServiceApiServer) CancelTransaction(ctx context.Context, request *gatewaypb.CancelTransactionRequest) (*gatewaypb.CancelTransactionResponse, error) {
	transactionId, err := DecodeTransactionId(request.TransactionId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid transaction id")
	}

	now := time.Now()

	resp1, err := s.coreApiClient.CancelTransaction(ctx, &corepb.CancelTransactionRequest{
		TransactionId: transactionId,
		Now:           now.UnixNano(),
	})
	if err != nil {
		return nil, monsterax.ErrorToGRPC(err)
	}

	return &gatewaypb.CancelTransactionResponse{
		Transaction: transactionToFront(resp1.Transaction),
	}, nil
}

func (s *LedgerServiceApiServer) CreateTransaction(ctx context.Context, request *gatewaypb.CreateTransactionRequest) (*gatewaypb.CreateTransactionResponse, error) {
	accountId, err := DecodeAccountId(request.AccountId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid account id")
	}

	now := time.Now()

	resp1, err := s.coreApiClient.CreateTransaction(ctx, &corepb.CreateTransactionRequest{
		TransactionId: &corepb.TransactionId{
			AccountId:     accountId,
			TransactionId: rand.Uint64(),
		},
		Amount:      request.Amount,
		Description: request.Description,
		Settled:     request.Settled,
		Now:         now.UnixNano(),
	})
	if err != nil {
		return nil, monsterax.ErrorToGRPC(err)
	}

	return &gatewaypb.CreateTransactionResponse{
		Transaction: transactionToFront(resp1.Transaction),
	}, nil
}

func (s *LedgerServiceApiServer) SettleTransaction(ctx context.Context, request *gatewaypb.SettleTransactionRequest) (*gatewaypb.SettleTransactionResponse, error) {
	transactionId, err := DecodeTransactionId(request.TransactionId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid transaction id")
	}

	now := time.Now()

	resp1, err := s.coreApiClient.SettleTransaction(ctx, &corepb.SettleTransactionRequest{
		TransactionId: transactionId,
		Now:           now.UnixNano(),
	})
	if err != nil {
		return nil, monsterax.ErrorToGRPC(err)
	}

	return &gatewaypb.SettleTransactionResponse{
		Transaction: transactionToFront(resp1.Transaction),
	}, nil
}

func NewLedgerServiceApiServer(coreApiClient LedgerServiceCoreApi) *LedgerServiceApiServer {
	return &LedgerServiceApiServer{
		coreApiClient: coreApiClient,
	}
}
