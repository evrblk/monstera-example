package ledger

import (
	"github.com/evrblk/monstera-example/ledger/corepb"
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

func TestTransactionIdEncodeDecode(t *testing.T) {
	require := require.New(t)

	for i := 0; i < 1000; i++ {
		id := &corepb.TransactionId{
			AccountId:     rand.Uint64(),
			TransactionId: rand.Uint64(),
		}

		actual, err := DecodeTransactionId(EncodeTransactionId(id))
		require.NoError(err)
		require.Equal(id, actual)
	}
}

func TestTransactionIdDecode(t *testing.T) {
	require := require.New(t)

	_, err := DecodeTransactionId("206af3de774e2Sc751d11665f8af8d41")
	require.Error(err)

	_, err = DecodeTransactionId("206af3de774e2dc751d11665f8af8d41")
	require.NoError(err)
}
