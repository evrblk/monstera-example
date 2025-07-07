package tinyurl

import (
	"github.com/evrblk/monstera-example/tinyurl/corepb"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserIdEncodeDecode(t *testing.T) {
	require := require.New(t)

	for i := 0; i < 1000; i++ {
		id := rand.Uint32()

		actual, err := DecodeUserId(EncodeUserId(id))
		require.NoError(err)
		require.Equal(id, actual)
	}
}

func TestUserIdDecode(t *testing.T) {
	require := require.New(t)

	_, err := DecodeUserId("c6ab9Xfb")
	require.Error(err)

	_, err = DecodeUserId("3d2dbe65")
	require.NoError(err)
}

func TestShortUrlIdEncodeDecode(t *testing.T) {
	require := require.New(t)

	for i := 0; i < 1000; i++ {
		id := &corepb.ShortUrlId{
			UserId:     rand.Uint32(),
			ShortUrlId: rand.Uint32(),
		}

		actual, err := DecodeShortUrlId(EncodeShortUrlId(id))
		require.NoError(err)
		require.Equal(id, actual)
	}
}

func TestShortUrlIdDecode(t *testing.T) {
	require := require.New(t)

	_, err := DecodeShortUrlId("206af3de774e2Sc7")
	require.Error(err)

	_, err = DecodeShortUrlId("206af3de774e2dc7")
	require.NoError(err)
}
