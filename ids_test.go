package monsteraexample

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccountIdEncodeDecode(t *testing.T) {
	require := require.New(t)

	for i := 0; i < 1000; i++ {
		id := rand.Uint64()

		actual, err := DecodeAccountId(EncodeAccountId(id))
		require.NoError(err)
		require.Equal(id, actual)
	}
}

func TestAccountIdDecode(t *testing.T) {
	require := require.New(t)

	_, err := DecodeAccountId("c6ab9cfb37117Sef")
	require.Error(err)

	_, err = DecodeAccountId("3d2dbe651ab536a1")
	require.NoError(err)
}
