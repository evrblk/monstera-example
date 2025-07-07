package tinyurl

import (
	"math/rand/v2"
	"testing"
	"time"

	"github.com/evrblk/monstera"
	"github.com/evrblk/monstera-example/tinyurl/corepb"
	"github.com/stretchr/testify/require"
)

func TestCreateAndGetUser(t *testing.T) {
	require := require.New(t)

	usersCore := newUsersCore()

	now := time.Now()

	userId := rand.Uint32()

	// create user
	response1, err := usersCore.CreateUser(&corepb.CreateUserRequest{
		UserId: userId,
		Now:    now.UnixNano(),
	})
	require.NoError(err)

	require.NotNil(response1.User)
	require.EqualValues(now.UnixNano(), response1.User.CreatedAt)
	require.EqualValues(now.UnixNano(), response1.User.UpdatedAt)
	require.Equal(userId, response1.User.Id)

	// get user
	response2, err := usersCore.GetUser(&corepb.GetUserRequest{
		UserId: userId,
	})
	require.NoError(err)
	require.NotNil(response2.User)
	require.EqualValues(now.UnixNano(), response2.User.CreatedAt)
	require.EqualValues(now.UnixNano(), response2.User.UpdatedAt)
	require.Equal(userId, response2.User.Id)
}

func newUsersCore() *UsersCore {
	return NewUsersCore(monstera.NewBadgerInMemoryStore(), []byte{0x00, 0x00}, []byte{0xff, 0xff})
}
