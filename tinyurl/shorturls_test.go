package tinyurl

import (
	"math/rand/v2"
	"testing"
	"time"

	"github.com/evrblk/monstera"
	"github.com/evrblk/monstera-example/tinyurl/corepb"
	"github.com/stretchr/testify/require"
)

func TestCreateAndGetShortUrl(t *testing.T) {
	require := require.New(t)

	shortUrlsCore := newShortUrlsCore()

	now := time.Now()

	userId := rand.Uint32()
	shortUrlId := rand.Uint32()

	// create short url
	response1, err := shortUrlsCore.CreateShortUrl(&corepb.CreateShortUrlRequest{
		ShortUrlId: &corepb.ShortUrlId{
			UserId:     userId,
			ShortUrlId: shortUrlId,
		},
		Now: now.UnixNano(),
	})
	require.NoError(err)

	require.NotNil(response1.ShortUrl)
	require.EqualValues(now.UnixNano(), response1.ShortUrl.CreatedAt)
	require.Equal(userId, response1.ShortUrl.Id.UserId)
	require.Equal(shortUrlId, response1.ShortUrl.Id.ShortUrlId)

	// get shortUrl
	response2, err := shortUrlsCore.GetShortUrl(&corepb.GetShortUrlRequest{
		ShortUrlId: &corepb.ShortUrlId{
			UserId:     userId,
			ShortUrlId: shortUrlId,
		},
	})
	require.NoError(err)
	require.NotNil(response2.ShortUrl)
	require.EqualValues(now.UnixNano(), response2.ShortUrl.CreatedAt)
	require.Equal(userId, response2.ShortUrl.Id.UserId)
	require.Equal(shortUrlId, response2.ShortUrl.Id.ShortUrlId)
}

func newShortUrlsCore() *ShortUrlsCore {
	return NewShortUrlsCore(monstera.NewBadgerInMemoryStore(), []byte{0x00, 0x00}, []byte{0xff, 0xff})
}
