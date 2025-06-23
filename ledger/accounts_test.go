package ledger

import (
	"math/rand/v2"
	"testing"
	"time"

	"github.com/evrblk/monstera"
	"github.com/evrblk/monstera-example/ledger/corepb"
	"github.com/stretchr/testify/require"
)

func TestCreateAndGetAccount(t *testing.T) {
	require := require.New(t)

	accountsCore := newAccountsCore()

	now := time.Now()

	accountId := rand.Uint64()

	// create account
	response1, err := accountsCore.CreateAccount(&corepb.CreateAccountRequest{
		AccountId: accountId,
		Now:       now.UnixNano(),
	})
	require.NoError(err)

	require.NotNil(response1.Account)
	require.EqualValues(now.UnixNano(), response1.Account.CreatedAt)
	require.EqualValues(now.UnixNano(), response1.Account.UpdatedAt)
	require.Equal(accountId, response1.Account.Id)
	require.EqualValues(0, response1.Account.AvailableBalance)
	require.EqualValues(0, response1.Account.SettledBalance)

	// get account
	response2, err := accountsCore.GetAccount(&corepb.GetAccountRequest{
		AccountId: accountId,
	})
	require.NoError(err)
	require.NotNil(response2.Account)
	require.EqualValues(now.UnixNano(), response2.Account.CreatedAt)
	require.EqualValues(now.UnixNano(), response2.Account.UpdatedAt)
	require.Equal(accountId, response2.Account.Id)
	require.EqualValues(0, response2.Account.AvailableBalance)
	require.EqualValues(0, response2.Account.SettledBalance)
}

func TestCreateAndSettleTransaction(t *testing.T) {
	require := require.New(t)

	accountsCore := newAccountsCore()

	now := time.Now()

	accountId := rand.Uint64()

	// T+0: create account
	_, err := accountsCore.CreateAccount(&corepb.CreateAccountRequest{
		AccountId: accountId,
		Now:       now.UnixNano(),
	})
	require.NoError(err)

	transactionId := rand.Uint64()

	// T+1m: create pending positive transaction
	response1, err := accountsCore.CreateTransaction(&corepb.CreateTransactionRequest{
		TransactionId: &corepb.TransactionId{
			AccountId:     accountId,
			TransactionId: transactionId,
		},
		Now:         now.Add(time.Minute).UnixNano(),
		Description: "Description",
		Amount:      100,
		Settled:     false,
	})
	require.NoError(err)
	require.NotNil(response1.Transaction)
	require.EqualValues(now.Add(time.Minute).UnixNano(), response1.Transaction.CreatedAt)
	require.EqualValues(now.Add(time.Minute).UnixNano(), response1.Transaction.UpdatedAt)
	require.Equal(corepb.TransactionStatus_TRANSACTION_STATUS_PENDING, response1.Transaction.Status)
	require.Equal(accountId, response1.Transaction.Id.AccountId)
	require.Equal(transactionId, response1.Transaction.Id.TransactionId)
	require.EqualValues(100, response1.Transaction.Amount)

	// T+1m: get account
	response2, err := accountsCore.GetAccount(&corepb.GetAccountRequest{
		AccountId: accountId,
	})
	require.NoError(err)
	require.NotNil(response2.Account)
	require.EqualValues(now.UnixNano(), response2.Account.CreatedAt)
	require.EqualValues(now.Add(time.Minute).UnixNano(), response2.Account.UpdatedAt)
	require.EqualValues(0, response2.Account.AvailableBalance)
	require.EqualValues(0, response2.Account.SettledBalance)

	// T+2m: settle transaction
	response3, err := accountsCore.SettleTransaction(&corepb.SettleTransactionRequest{
		TransactionId: response1.Transaction.Id,
		Now:           now.Add(2 * time.Minute).UnixNano(),
	})
	require.NoError(err)
	require.NotNil(response3.Transaction)
	require.EqualValues(now.Add(time.Minute).UnixNano(), response3.Transaction.CreatedAt)
	require.EqualValues(now.Add(2*time.Minute).UnixNano(), response3.Transaction.UpdatedAt)
	require.Equal(corepb.TransactionStatus_TRANSACTION_STATUS_SETTLED, response3.Transaction.Status)

	// T+2m: get account
	response4, err := accountsCore.GetAccount(&corepb.GetAccountRequest{
		AccountId: accountId,
	})
	require.NoError(err)
	require.NotNil(response4.Account)
	require.EqualValues(now.UnixNano(), response4.Account.CreatedAt)
	require.EqualValues(now.Add(2*time.Minute).UnixNano(), response4.Account.UpdatedAt)
	require.EqualValues(100, response4.Account.AvailableBalance)
	require.EqualValues(100, response4.Account.SettledBalance)

	// T+3m: create pending negative transaction
	response5, err := accountsCore.CreateTransaction(&corepb.CreateTransactionRequest{
		TransactionId: &corepb.TransactionId{
			AccountId:     accountId,
			TransactionId: transactionId,
		},
		Now:         now.Add(3 * time.Minute).UnixNano(),
		Description: "Description",
		Amount:      -10,
		Settled:     false,
	})
	require.NoError(err)
	require.NotNil(response5.Transaction)
	require.Equal(corepb.TransactionStatus_TRANSACTION_STATUS_PENDING, response5.Transaction.Status)
	require.EqualValues(-10, response5.Transaction.Amount)

	// T+3m: get account
	response6, err := accountsCore.GetAccount(&corepb.GetAccountRequest{
		AccountId: accountId,
	})
	require.NoError(err)
	require.NotNil(response6.Account)
	require.EqualValues(now.UnixNano(), response6.Account.CreatedAt)
	require.EqualValues(now.Add(3*time.Minute).UnixNano(), response6.Account.UpdatedAt)
	require.EqualValues(90, response6.Account.AvailableBalance)
	require.EqualValues(100, response6.Account.SettledBalance)

	// T+4m: settle transaction
	response7, err := accountsCore.SettleTransaction(&corepb.SettleTransactionRequest{
		TransactionId: response5.Transaction.Id,
		Now:           now.Add(4 * time.Minute).UnixNano(),
	})
	require.NoError(err)
	require.NotNil(response7.Transaction)
	require.Equal(corepb.TransactionStatus_TRANSACTION_STATUS_SETTLED, response7.Transaction.Status)

	// T+4m: get account
	response8, err := accountsCore.GetAccount(&corepb.GetAccountRequest{
		AccountId: accountId,
	})
	require.NoError(err)
	require.NotNil(response8.Account)
	require.EqualValues(now.UnixNano(), response8.Account.CreatedAt)
	require.EqualValues(now.Add(4*time.Minute).UnixNano(), response8.Account.UpdatedAt)
	require.EqualValues(90, response8.Account.AvailableBalance)
	require.EqualValues(90, response8.Account.SettledBalance)
}

func newAccountsCore() *AccountsCore {
	return NewAccountsCore(monstera.NewBadgerInMemoryStore(), []byte{0x00, 0x00}, []byte{0xff, 0xff})
}
