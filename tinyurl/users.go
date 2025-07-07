package tinyurl

import (
	"io"

	"errors"

	"github.com/evrblk/monstera"
	"github.com/evrblk/monstera-example/tinyurl/corepb"
	monsterax "github.com/evrblk/monstera/x"
)

type UsersCore struct {
	badgerStore *monstera.BadgerStore

	usersTable *monsterax.SimpleKeyTable[*corepb.User, corepb.User]
}

var _ UsersCoreApi = &UsersCore{}

func NewUsersCore(badgerStore *monstera.BadgerStore, shardLowerBound []byte, shardUpperBound []byte) *UsersCore {
	return &UsersCore{
		badgerStore: badgerStore,
		usersTable:  monsterax.NewSimpleKeyTable[*corepb.User, corepb.User](usersTableId, shardLowerBound, shardUpperBound),
	}
}

func (c *UsersCore) ranges() []monstera.KeyRange {
	return []monstera.KeyRange{
		c.usersTable.GetTableKeyRange(),
	}
}

func (c *UsersCore) Snapshot() monstera.ApplicationCoreSnapshot {
	return monsterax.Snapshot(c.badgerStore, c.ranges())
}

func (c *UsersCore) Restore(reader io.ReadCloser) error {
	return monsterax.Restore(c.badgerStore, c.ranges(), reader)
}

func (c *UsersCore) Close() {

}

func (c *UsersCore) GetUser(request *corepb.GetUserRequest) (*corepb.GetUserResponse, error) {
	txn := c.badgerStore.View()
	defer txn.Discard()

	user, err := c.getUser(txn, request.UserId)
	if err != nil {
		if errors.Is(err, monstera.ErrNotFound) {
			return nil, monsterax.NewErrorWithContext(
				monsterax.NotFound,
				"user not found",
				map[string]string{"user_id": EncodeUserId(request.UserId)})
		} else {
			panic(err)
		}
	}

	return &corepb.GetUserResponse{
		User: user,
	}, nil
}

func (c *UsersCore) CreateUser(request *corepb.CreateUserRequest) (*corepb.CreateUserResponse, error) {
	txn := c.badgerStore.Update()
	defer txn.Discard()

	user := &corepb.User{
		Id:        request.UserId,
		CreatedAt: request.Now,
		UpdatedAt: request.Now,
	}

	err := c.createUser(txn, user)
	panicIfNotNil(err)

	err = txn.Commit()
	panicIfNotNil(err)

	return &corepb.CreateUserResponse{
		User: user,
	}, nil
}

func (c *UsersCore) getUser(txn *monstera.Txn, userId uint32) (*corepb.User, error) {
	return c.usersTable.Get(txn, usersTablePK(userId))
}

func (c *UsersCore) updateUser(txn *monstera.Txn, user *corepb.User) error {
	return c.usersTable.Set(txn, usersTablePK(user.Id), user)
}

func (c *UsersCore) deleteUser(txn *monstera.Txn, userId uint32) error {
	return c.usersTable.Delete(txn, usersTablePK(userId))
}

func (c *UsersCore) createUser(txn *monstera.Txn, user *corepb.User) error {
	return c.usersTable.Set(txn, usersTablePK(user.Id), user)
}

// 1. shard key (by user id)
// 2. user id
func usersTablePK(userId uint32) []byte {
	return monstera.ConcatBytes(shardByUser(userId), userId)
}

func panicIfNotNil(err error) {
	if err != nil {
		panic(err)
	}
}
