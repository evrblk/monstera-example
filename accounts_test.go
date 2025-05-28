package monsteraexample

import (
	"math/rand"
	"testing"
	"time"

	"github.com/evrblk/monstera"
	"github.com/evrblk/monstera-example/corepb"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestCreateAndGetAccount(t *testing.T) {
	require := require.New(t)

	accountsCore := newAccountsCore()

	now := time.Now()
	accountId := rand.Uint64()

	// Create account
	response1, err := accountsCore.CreateAccount(&corepb.CreateAccountRequest{
		AccountId: accountId,
		FullName:  "Doogie Howser",
		Email:     "doogie@gmail.com",
		Now:       now.UnixNano(),
	})

	require.NoError(err)
	require.NotNil(response1.Account)

	// Get this newly created account
	response2, err := accountsCore.GetAccount(&corepb.GetAccountRequest{
		AccountId: accountId,
	})

	require.NoError(err)
	require.NotNil(response2.Account)

	require.Equal(accountId, response2.Account.Id)
	require.Equal("doogie@gmail.com", response2.Account.Email)
	require.Equal(now.UnixNano(), response2.Account.CreatedAt)
	require.Equal(now.UnixNano(), response2.Account.UpdatedAt)
	require.Equal("Doogie Howser", response2.Account.FullName)

	// Get non-existent account
	_, err = accountsCore.GetAccount(&corepb.GetAccountRequest{
		AccountId: rand.Uint64(),
	})

	require.Error(err)
}

func TestExistingEmail(t *testing.T) {
	t.SkipNow()

	require := require.New(t)

	accountsCore := newAccountsCore()

	now := time.Now()

	// Create account 1
	response1, err := accountsCore.CreateAccount(&corepb.CreateAccountRequest{
		AccountId: rand.Uint64(),
		FullName:  "Squirrel McFuzzy",
		Email:     "squirrel@gmail.com",
		Now:       now.UnixNano(),
	})

	require.NoError(err)
	require.NotNil(response1.Account)

	// Create account 2 with another email
	response2, err := accountsCore.CreateAccount(&corepb.CreateAccountRequest{
		AccountId: rand.Uint64(),
		FullName:  "Doogie Howser",
		Email:     "doogie@gmail.com",
		Now:       now.Add(time.Second).UnixNano(),
	})

	require.NoError(err)
	require.NotNil(response2.Account)

	// Create account 3 with the same email
	_, err = accountsCore.CreateAccount(&corepb.CreateAccountRequest{
		AccountId: rand.Uint64(),
		FullName:  "Doogie Howser Jr.",
		Email:     "doogie@gmail.com",
		Now:       now.Add(time.Second * 2).UnixNano(),
	})

	require.Error(err)
}

func TestListAccounts(t *testing.T) {
	require := require.New(t)

	accountsCore := newAccountsCore()

	now := time.Now()

	// Create account 1
	response1, err := accountsCore.CreateAccount(&corepb.CreateAccountRequest{
		AccountId: rand.Uint64(),
		FullName:  "Squirrel McFuzzy",
		Email:     "squirrel@gmail.com",
		Now:       now.UnixNano(),
	})

	require.NoError(err)

	// Create account 2 with another email
	response2, err := accountsCore.CreateAccount(&corepb.CreateAccountRequest{
		AccountId: rand.Uint64(),
		FullName:  "Doogie Howser",
		Email:     "doogie@gmail.com",
		Now:       now.Add(time.Second).UnixNano(),
	})

	require.NoError(err)

	// List accountsCore
	response3, err := accountsCore.ListAccounts(&corepb.ListAccountsRequest{})

	require.NoError(err)
	require.Len(response3.Accounts, 2)
	ids := lo.Map(response3.Accounts, func(a *corepb.Account, index int) uint64 {
		return a.Id
	})
	require.ElementsMatch(ids, []uint64{response1.Account.Id, response2.Account.Id})
}

func newAccountsCore() *AccountsCore {
	return NewAccountsCore(monstera.NewBadgerInMemoryStore())
}
