// Code generated by `monstera generate`. DO NOT EDIT.

package ledger

import (
	monstera "github.com/evrblk/monstera"
	corepb "github.com/evrblk/monstera-example/ledger/corepb"
	monsterax "github.com/evrblk/monstera/x"
	proto "google.golang.org/protobuf/proto"
	"io"
)

type AccountsCoreAdapter struct {
	accountsCore AccountsCoreApi
}

var _ monstera.ApplicationCore = &AccountsCoreAdapter{}

func NewAccountsCoreAdapter(accountsCore AccountsCoreApi) *AccountsCoreAdapter {
	return &AccountsCoreAdapter{accountsCore: accountsCore}
}

func (a *AccountsCoreAdapter) Snapshot() monstera.ApplicationCoreSnapshot {
	return a.accountsCore.Snapshot()
}

func (a *AccountsCoreAdapter) Restore(r io.ReadCloser) error {
	return a.accountsCore.Restore(r)
}

func (a *AccountsCoreAdapter) Close() {
	a.accountsCore.Close()
}

func (a *AccountsCoreAdapter) Update(request []byte) []byte {
	updateRequest := &corepb.UpdateRequest{}
	updateResponse := &corepb.UpdateResponse{}

	err := proto.Unmarshal(request, updateRequest)
	if err != nil {
		panic(err)
	}

	switch req := updateRequest.Request.(type) {
	case *corepb.UpdateRequest_CreateTransactionRequest:
		r, err := a.accountsCore.CreateTransaction(req.CreateTransactionRequest)
		updateResponse.Response = &corepb.UpdateResponse_CreateTransactionResponse{CreateTransactionResponse: r}
		updateResponse.Error = monsterax.WrapError(err)
	case *corepb.UpdateRequest_CancelTransactionRequest:
		r, err := a.accountsCore.CancelTransaction(req.CancelTransactionRequest)
		updateResponse.Response = &corepb.UpdateResponse_CancelTransactionResponse{CancelTransactionResponse: r}
		updateResponse.Error = monsterax.WrapError(err)
	case *corepb.UpdateRequest_SettleTransactionRequest:
		r, err := a.accountsCore.SettleTransaction(req.SettleTransactionRequest)
		updateResponse.Response = &corepb.UpdateResponse_SettleTransactionResponse{SettleTransactionResponse: r}
		updateResponse.Error = monsterax.WrapError(err)
	case *corepb.UpdateRequest_CreateAccountRequest:
		r, err := a.accountsCore.CreateAccount(req.CreateAccountRequest)
		updateResponse.Response = &corepb.UpdateResponse_CreateAccountResponse{CreateAccountResponse: r}
		updateResponse.Error = monsterax.WrapError(err)
	default:
		panic("no matching handlers")
	}
	response, err := proto.Marshal(updateResponse)
	if err != nil {
		panic(err)
	}

	return response
}

func (a *AccountsCoreAdapter) Read(request []byte) []byte {
	readRequest := &corepb.ReadRequest{}
	readResponse := &corepb.ReadResponse{}

	err := proto.Unmarshal(request, readRequest)
	if err != nil {
		panic(err)
	}

	switch req := readRequest.Request.(type) {
	case *corepb.ReadRequest_ListTransactionsRequest:
		r, err := a.accountsCore.ListTransactions(req.ListTransactionsRequest)
		readResponse.Response = &corepb.ReadResponse_ListTransactionsResponse{ListTransactionsResponse: r}
		readResponse.Error = monsterax.WrapError(err)
	case *corepb.ReadRequest_GetTransactionRequest:
		r, err := a.accountsCore.GetTransaction(req.GetTransactionRequest)
		readResponse.Response = &corepb.ReadResponse_GetTransactionResponse{GetTransactionResponse: r}
		readResponse.Error = monsterax.WrapError(err)
	case *corepb.ReadRequest_GetAccountRequest:
		r, err := a.accountsCore.GetAccount(req.GetAccountRequest)
		readResponse.Response = &corepb.ReadResponse_GetAccountResponse{GetAccountResponse: r}
		readResponse.Error = monsterax.WrapError(err)
	default:
		panic("no matching handlers")
	}
	response, err := proto.Marshal(readResponse)
	if err != nil {
		panic(err)
	}

	return response
}
