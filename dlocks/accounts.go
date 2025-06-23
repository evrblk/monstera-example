package dlocks

import (
	"encoding/binary"
	"errors"
	"io"
	"strings"

	"github.com/evrblk/monstera"
	"github.com/evrblk/monstera-example/dlocks/corepb"
	monsterax "github.com/evrblk/monstera/x"
)

type AccountsCore struct {
	badgerStore *monstera.BadgerStore

	accountsTable       *monsterax.SimpleKeyTable[*corepb.Account, corepb.Account]
	accountsEmailsIndex *monsterax.UniqueUint64Index
}

func NewAccountsCore(badgerStore *monstera.BadgerStore) *AccountsCore {
	return &AccountsCore{
		badgerStore: badgerStore,

		accountsTable:       monsterax.NewSimpleKeyTable[*corepb.Account, corepb.Account](accountsTableId, []byte{0x00}, []byte{0xff}),
		accountsEmailsIndex: monsterax.NewUniqueUint64Index(accountsEmailsIndexId, []byte{0x00}, []byte{0xff}),
	}
}

func (c *AccountsCore) ranges() []monstera.KeyRange {
	return []monstera.KeyRange{
		c.accountsTable.GetTableKeyRange(),
		c.accountsEmailsIndex.GetTableKeyRange(),
	}
}

func (c *AccountsCore) Close() {

}

func (c *AccountsCore) Restore(reader io.ReadCloser) error {
	return monsterax.Restore(c.badgerStore, c.ranges(), reader)
}

func (c *AccountsCore) Snapshot() monstera.ApplicationCoreSnapshot {
	return monsterax.Snapshot(c.badgerStore, c.ranges())
}

func (c *AccountsCore) CreateAccount(request *corepb.CreateAccountRequest) (*corepb.CreateAccountResponse, error) {
	txn := c.badgerStore.Update()
	defer txn.Discard()

	// Checking email uniqueness
	email := strings.TrimSpace(strings.ToLower(request.Email))
	_, err := c.getAccountByEmail(txn, email)
	if err == nil {
		// Account with such email already exists
		return nil, monsterax.NewErrorWithContext(monsterax.AlreadyExists, "account with this email already exists", map[string]string{"email": request.Email})
	} else if !errors.Is(err, monstera.ErrNotFound) {
		panic(err)
	}

	account := &corepb.Account{
		Id:                    request.AccountId,
		Email:                 email,
		FullName:              request.FullName,
		CreatedAt:             request.Now,
		UpdatedAt:             request.Now,
		MaxNumberOfNamespaces: request.MaxNumberOfNamespaces,
	}

	err = c.createAccount(txn, account)
	panicIfNotNil(err)

	err = txn.Commit()
	panicIfNotNil(err)

	return &corepb.CreateAccountResponse{
		Account: account,
	}, nil
}

func (c *AccountsCore) ListAccounts(request *corepb.ListAccountsRequest) (*corepb.ListAccountsResponse, error) {
	txn := c.badgerStore.View()
	defer txn.Discard()

	accounts, err := c.listAccounts(txn)
	panicIfNotNil(err)

	return &corepb.ListAccountsResponse{
		Accounts: accounts,
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
				map[string]string{"accounts_id": EncodeAccountId(request.AccountId)})
		} else {
			panic(err)
		}
	}

	return &corepb.GetAccountResponse{
		Account: account,
	}, nil
}

func (c *AccountsCore) UpdateAccount(request *corepb.UpdateAccountRequest) (*corepb.UpdateAccountResponse, error) {
	txn := c.badgerStore.Update()
	defer txn.Discard()

	account, err := c.getAccount(txn, request.AccountId)
	if err != nil {
		if errors.Is(err, monstera.ErrNotFound) {
			return nil, monsterax.NewErrorWithContext(
				monsterax.NotFound,
				"account not found",
				map[string]string{"accounts_id": EncodeAccountId(request.AccountId)})
		} else {
			panic(err)
		}
	}

	account.UpdatedAt = request.Now
	account.FullName = request.FullName

	err = c.updateAccount(txn, account)
	panicIfNotNil(err)

	err = txn.Commit()
	panicIfNotNil(err)

	return &corepb.UpdateAccountResponse{
		Account: account,
	}, nil
}

func (c *AccountsCore) DeleteAccount(request *corepb.DeleteAccountRequest) (*corepb.DeleteAccountResponse, error) {
	txn := c.badgerStore.Update()
	defer txn.Discard()

	account, err := c.getAccount(txn, request.AccountId)
	if err != nil {
		if errors.Is(err, monstera.ErrNotFound) {
			return nil, monsterax.NewErrorWithContext(
				monsterax.NotFound,
				"account not found",
				map[string]string{"accounts_id": EncodeAccountId(request.AccountId)})
		} else {
			panic(err)
		}
	}

	err = c.deleteAccount(txn, account)
	panicIfNotNil(err)

	err = txn.Commit()
	panicIfNotNil(err)

	return &corepb.DeleteAccountResponse{}, nil
}

func (c *AccountsCore) getAccount(txn *monstera.Txn, accountId uint64) (*corepb.Account, error) {
	return c.accountsTable.Get(txn, accountIdAsPK(accountId))
}

func (c *AccountsCore) getAccountByEmail(txn *monstera.Txn, email string) (*corepb.Account, error) {
	accountId, err := c.accountsEmailsIndex.Get(txn, []byte(email))
	if err != nil {
		return nil, err
	}
	return c.getAccount(txn, accountId)
}

func (c *AccountsCore) createAccount(txn *monstera.Txn, account *corepb.Account) error {
	err := c.accountsEmailsIndex.Set(txn, []byte(account.Email), account.Id)
	if err != nil {
		return err
	}

	return c.accountsTable.Set(txn, accountIdAsPK(account.Id), account)
}

func (c *AccountsCore) updateAccount(txn *monstera.Txn, account *corepb.Account) error {
	return c.accountsTable.Set(txn, accountIdAsPK(account.Id), account)
}

func (c *AccountsCore) deleteAccount(txn *monstera.Txn, account *corepb.Account) error {
	return c.accountsTable.Delete(txn, accountIdAsPK(account.Id))
}

func (c *AccountsCore) listAccounts(txn *monstera.Txn) ([]*corepb.Account, error) {
	result := make([]*corepb.Account, 0)

	err := c.accountsTable.List(txn, func(account *corepb.Account) (bool, error) {
		result = append(result, account)
		return true, nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func accountIdAsPK(accountId uint64) []byte {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key[0:], accountId)
	return key
}
