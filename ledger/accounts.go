package ledger

import (
	"io"

	"errors"

	"github.com/evrblk/monstera"
	"github.com/evrblk/monstera-example/ledger/corepb"
	monsterax "github.com/evrblk/monstera/x"
)

type AccountsCore struct {
	badgerStore *monstera.BadgerStore

	accountsTable     *monsterax.SimpleKeyTable[*corepb.Account, corepb.Account]
	transactionsTable *monsterax.CompositeKeyTable[*corepb.Transaction, corepb.Transaction]
}

var _ AccountsCoreApi = &AccountsCore{}

func NewAccountsCore(badgerStore *monstera.BadgerStore, shardLowerBound []byte, shardUpperBound []byte) *AccountsCore {
	return &AccountsCore{
		badgerStore:       badgerStore,
		accountsTable:     monsterax.NewSimpleKeyTable[*corepb.Account, corepb.Account](accountsTableId, shardLowerBound, shardUpperBound),
		transactionsTable: monsterax.NewCompositeKeyTable[*corepb.Transaction, corepb.Transaction](transactionsTableId, shardLowerBound, shardUpperBound),
	}
}

func (c *AccountsCore) ranges() []monstera.KeyRange {
	return []monstera.KeyRange{
		c.accountsTable.GetTableKeyRange(),
		c.transactionsTable.GetTableKeyRange(),
	}
}

func (c *AccountsCore) Snapshot() monstera.ApplicationCoreSnapshot {
	return monsterax.Snapshot(c.badgerStore, c.ranges())
}

func (c *AccountsCore) Restore(reader io.ReadCloser) error {
	return monsterax.Restore(c.badgerStore, c.ranges(), reader)
}

func (c *AccountsCore) Close() {

}

func (c *AccountsCore) ListTransactions(request *corepb.ListTransactionsRequest) (*corepb.ListTransactionsResponse, error) {
	txn := c.badgerStore.View()
	defer txn.Discard()

	transactions, err := c.listTransactions(txn, request.AccountId)
	panicIfNotNil(err)

	return &corepb.ListTransactionsResponse{
		Transactions: transactions,
	}, nil
}

func (c *AccountsCore) GetTransaction(request *corepb.GetTransactionRequest) (*corepb.GetTransactionResponse, error) {
	txn := c.badgerStore.View()
	defer txn.Discard()

	transaction, err := c.getTransaction(txn, request.TransactionId)
	if err != nil {
		if errors.Is(err, monstera.ErrNotFound) {
			return nil, monsterax.NewErrorWithContext(
				monsterax.NotFound,
				"transaction not found",
				map[string]string{"transaction_id": EncodeTransactionId(request.TransactionId)})
		} else {
			panic(err)
		}
	}

	return &corepb.GetTransactionResponse{
		Transaction: transaction,
	}, nil
}

func (c *AccountsCore) CreateTransaction(request *corepb.CreateTransactionRequest) (*corepb.CreateTransactionResponse, error) {
	txn := c.badgerStore.Update()
	defer txn.Discard()

	account, err := c.getAccount(txn, request.TransactionId.AccountId)
	if err != nil {
		if errors.Is(err, monstera.ErrNotFound) {
			return nil, monsterax.NewErrorWithContext(
				monsterax.NotFound,
				"account not found",
				map[string]string{"account_id": EncodeAccountId(request.TransactionId.AccountId)})
		} else {
			panic(err)
		}
	}

	transaction := &corepb.Transaction{
		Id:          request.TransactionId,
		Amount:      request.Amount,
		Description: request.Description,
		CreatedAt:   request.Now,
		UpdatedAt:   request.Now,
	}

	// if money is added to the account (debits)
	if transaction.Amount >= 0 {
		if request.Settled {
			transaction.Status = corepb.TransactionStatus_TRANSACTION_STATUS_SETTLED

			// settled debits increase the balance
			account.AvailableBalance += transaction.Amount
			account.SettledBalance += transaction.Amount
		} else {
			transaction.Status = corepb.TransactionStatus_TRANSACTION_STATUS_PENDING
			// pending debits do not increase the balance
		}
	} else {
		if account.AvailableBalance+transaction.Amount < 0 {
			// negative balance is not allowed for credits
			transaction.Status = corepb.TransactionStatus_TRANSACTION_STATUS_INSUFFICIENT_FUNDS
		} else {
			if request.Settled {
				transaction.Status = corepb.TransactionStatus_TRANSACTION_STATUS_SETTLED

				// settled credits decrease the balance (amount is negative)
				account.AvailableBalance += transaction.Amount
				account.SettledBalance += transaction.Amount
			} else {
				transaction.Status = corepb.TransactionStatus_TRANSACTION_STATUS_PENDING

				// pending credits decrease available balance only
				account.AvailableBalance += transaction.Amount
			}
		}
	}

	account.UpdatedAt = request.Now

	err = c.updateAccount(txn, account)
	panicIfNotNil(err)

	err = c.createTransaction(txn, transaction)
	panicIfNotNil(err)

	err = txn.Commit()
	panicIfNotNil(err)

	return &corepb.CreateTransactionResponse{
		Transaction: transaction,
	}, nil
}

func (c *AccountsCore) CancelTransaction(request *corepb.CancelTransactionRequest) (*corepb.CancelTransactionResponse, error) {
	txn := c.badgerStore.Update()
	defer txn.Discard()

	account, err := c.getAccount(txn, request.TransactionId.AccountId)
	if err != nil {
		if errors.Is(err, monstera.ErrNotFound) {
			return nil, monsterax.NewErrorWithContext(
				monsterax.NotFound,
				"account not found",
				map[string]string{"account_id": EncodeAccountId(request.TransactionId.AccountId)})
		} else {
			panic(err)
		}
	}

	transaction, err := c.getTransaction(txn, request.TransactionId)
	if err != nil {
		if errors.Is(err, monstera.ErrNotFound) {
			return nil, monsterax.NewErrorWithContext(
				monsterax.NotFound,
				"transaction not found",
				map[string]string{"transaction_id": EncodeTransactionId(request.TransactionId)})
		} else {
			panic(err)
		}
	}

	// only pending transactions can be canceled
	if transaction.Status != corepb.TransactionStatus_TRANSACTION_STATUS_PENDING {
		return nil, monsterax.NewErrorWithContext(
			monsterax.NotFound,
			"transaction is not pending",
			map[string]string{"transaction_id": EncodeTransactionId(request.TransactionId)})
	}

	transaction.Status = corepb.TransactionStatus_TRANSACTION_STATUS_CANCELLED
	transaction.UpdatedAt = request.Now

	// pending debits do not change the balance, no need to do anything
	// pending credits only change the available balance, need to subtract transaction amount
	if transaction.Amount < 0 {
		account.AvailableBalance -= transaction.Amount
	}

	account.UpdatedAt = request.Now

	err = c.updateTransaction(txn, transaction)
	panicIfNotNil(err)

	err = c.updateAccount(txn, account)
	panicIfNotNil(err)

	err = txn.Commit()
	panicIfNotNil(err)

	return &corepb.CancelTransactionResponse{
		Transaction: transaction,
	}, nil
}

func (c *AccountsCore) SettleTransaction(request *corepb.SettleTransactionRequest) (*corepb.SettleTransactionResponse, error) {
	txn := c.badgerStore.Update()
	defer txn.Discard()

	account, err := c.getAccount(txn, request.TransactionId.AccountId)
	if err != nil {
		if errors.Is(err, monstera.ErrNotFound) {
			return nil, monsterax.NewErrorWithContext(
				monsterax.NotFound,
				"account not found",
				map[string]string{"account_id": EncodeAccountId(request.TransactionId.AccountId)})
		} else {
			panic(err)
		}
	}

	transaction, err := c.getTransaction(txn, request.TransactionId)
	if err != nil {
		if errors.Is(err, monstera.ErrNotFound) {
			return nil, monsterax.NewErrorWithContext(
				monsterax.NotFound,
				"transaction not found",
				map[string]string{"transaction_id": EncodeTransactionId(request.TransactionId)})
		} else {
			panic(err)
		}
	}

	// only pending transactions can be settled
	if transaction.Status != corepb.TransactionStatus_TRANSACTION_STATUS_PENDING {
		return nil, monsterax.NewErrorWithContext(
			monsterax.NotFound,
			"transaction is not pending",
			map[string]string{"transaction_id": EncodeTransactionId(request.TransactionId)})
	}

	transaction.Status = corepb.TransactionStatus_TRANSACTION_STATUS_SETTLED
	transaction.UpdatedAt = request.Now

	// pending debits do not change the balance, need to update both balances
	// pending credits only change the available balance, need to update settled balance
	if transaction.Amount >= 0 {
		account.AvailableBalance += transaction.Amount
		account.SettledBalance += transaction.Amount
	} else {
		account.SettledBalance += transaction.Amount
	}

	account.UpdatedAt = request.Now

	err = c.updateTransaction(txn, transaction)
	panicIfNotNil(err)

	err = c.updateAccount(txn, account)
	panicIfNotNil(err)

	err = txn.Commit()
	panicIfNotNil(err)

	return &corepb.SettleTransactionResponse{
		Transaction: transaction,
	}, nil
}

func (c *AccountsCore) GetAccount(request *corepb.GetAccountRequest) (*corepb.GetAccountResponse, error) {
	txn := c.badgerStore.View()
	defer txn.Discard()

	account, err := c.getAccount(txn, request.AccountId)
	if err != nil {
		if errors.Is(err, monstera.ErrNotFound) {
			return nil, monsterax.NewErrorWithContext(
				monsterax.NotFound,
				"account not found",
				map[string]string{"account_id": EncodeAccountId(request.AccountId)})
		} else {
			panic(err)
		}
	}

	return &corepb.GetAccountResponse{
		Account: account,
	}, nil
}

func (c *AccountsCore) CreateAccount(request *corepb.CreateAccountRequest) (*corepb.CreateAccountResponse, error) {
	txn := c.badgerStore.Update()
	defer txn.Discard()

	account := &corepb.Account{
		Id:               request.AccountId,
		AvailableBalance: 0,
		SettledBalance:   0,
		CreatedAt:        request.Now,
		UpdatedAt:        request.Now,
	}

	err := c.createAccount(txn, account)
	panicIfNotNil(err)

	err = txn.Commit()
	panicIfNotNil(err)

	return &corepb.CreateAccountResponse{
		Account: account,
	}, nil
}

func (c *AccountsCore) getAccount(txn *monstera.Txn, accountId uint64) (*corepb.Account, error) {
	return c.accountsTable.Get(txn, accountsTablePK(accountId))
}

func (c *AccountsCore) updateAccount(txn *monstera.Txn, account *corepb.Account) error {
	return c.accountsTable.Set(txn, accountsTablePK(account.Id), account)
}

func (c *AccountsCore) deleteAccount(txn *monstera.Txn, accountId uint64) error {
	return c.accountsTable.Delete(txn, accountsTablePK(accountId))
}

func (c *AccountsCore) createAccount(txn *monstera.Txn, account *corepb.Account) error {
	return c.accountsTable.Set(txn, accountsTablePK(account.Id), account)
}

func (c *AccountsCore) getTransaction(txn *monstera.Txn, transactionId *corepb.TransactionId) (*corepb.Transaction, error) {
	return c.transactionsTable.Get(txn, transactionsTablePK(transactionId.AccountId), transactionsTableSK(transactionId))
}

func (c *AccountsCore) createTransaction(txn *monstera.Txn, transaction *corepb.Transaction) error {
	return c.transactionsTable.Set(txn, transactionsTablePK(transaction.Id.AccountId), transactionsTableSK(transaction.Id), transaction)
}

func (c *AccountsCore) updateTransaction(txn *monstera.Txn, transaction *corepb.Transaction) error {
	return c.transactionsTable.Set(txn, transactionsTablePK(transaction.Id.AccountId), transactionsTableSK(transaction.Id), transaction)
}

func (c *AccountsCore) listTransactions(txn *monstera.Txn, accountId uint64) ([]*corepb.Transaction, error) {
	result := make([]*corepb.Transaction, 0)

	err := c.transactionsTable.List(txn, transactionsTablePK(accountId), func(transaction *corepb.Transaction) (bool, error) {
		result = append(result, transaction)
		return true, nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// 1. shard key (by account id)
// 2. account id
func accountsTablePK(accountId uint64) []byte {
	return monstera.ConcatBytes(shardByAccount(accountId), accountId)
}

// 1. shard key (by account id)
// 2. account id
func transactionsTablePK(accountId uint64) []byte {
	return monstera.ConcatBytes(shardByAccount(accountId), accountId)
}

// 1. transaction id
func transactionsTableSK(t *corepb.TransactionId) []byte {
	return monstera.ConcatBytes(t.GetTransactionId())
}

func panicIfNotNil(err error) {
	if err != nil {
		panic(err)
	}
}
