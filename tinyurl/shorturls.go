package tinyurl

import (
	"io"

	"errors"

	"github.com/evrblk/monstera"
	"github.com/evrblk/monstera-example/tinyurl/corepb"
	monsterax "github.com/evrblk/monstera/x"
)

type ShortUrlsCore struct {
	badgerStore *monstera.BadgerStore

	shortUrlsTable *monsterax.CompositeKeyTable[*corepb.ShortUrl, corepb.ShortUrl]
}

var _ ShortUrlsCoreApi = &ShortUrlsCore{}

func NewShortUrlsCore(badgerStore *monstera.BadgerStore, shardLowerBound []byte, shardUpperBound []byte) *ShortUrlsCore {
	return &ShortUrlsCore{
		badgerStore:    badgerStore,
		shortUrlsTable: monsterax.NewCompositeKeyTable[*corepb.ShortUrl, corepb.ShortUrl](usersTableId, shardLowerBound, shardUpperBound),
	}
}

func (c *ShortUrlsCore) ranges() []monstera.KeyRange {
	return []monstera.KeyRange{
		c.shortUrlsTable.GetTableKeyRange(),
	}
}

func (c *ShortUrlsCore) Snapshot() monstera.ApplicationCoreSnapshot {
	return monsterax.Snapshot(c.badgerStore, c.ranges())
}

func (c *ShortUrlsCore) Restore(reader io.ReadCloser) error {
	return monsterax.Restore(c.badgerStore, c.ranges(), reader)
}

func (c *ShortUrlsCore) Close() {

}

func (c *ShortUrlsCore) GetShortUrl(request *corepb.GetShortUrlRequest) (*corepb.GetShortUrlResponse, error) {
	txn := c.badgerStore.View()
	defer txn.Discard()

	shortUrl, err := c.getShortUrl(txn, request.ShortUrlId)
	if err != nil {
		if errors.Is(err, monstera.ErrNotFound) {
			return nil, monsterax.NewErrorWithContext(
				monsterax.NotFound,
				"shortUrl not found",
				map[string]string{"short_url_id": EncodeShortUrlId(request.ShortUrlId)})
		} else {
			panic(err)
		}
	}

	return &corepb.GetShortUrlResponse{
		ShortUrl: shortUrl,
	}, nil
}

func (c *ShortUrlsCore) ListShortUrls(request *corepb.ListShortUrlsRequest) (*corepb.ListShortUrlsResponse, error) {
	txn := c.badgerStore.View()
	defer txn.Discard()

	shortUrls, err := c.listShortUrls(txn, request.UserId)
	panicIfNotNil(err)

	return &corepb.ListShortUrlsResponse{
		ShortUrls: shortUrls,
	}, nil
}

func (c *ShortUrlsCore) CreateShortUrl(request *corepb.CreateShortUrlRequest) (*corepb.CreateShortUrlResponse, error) {
	txn := c.badgerStore.Update()
	defer txn.Discard()

	shortUrl := &corepb.ShortUrl{
		Id:        request.ShortUrlId,
		FullUrl:   request.FullUrl,
		CreatedAt: request.Now,
	}

	err := c.createShortUrl(txn, shortUrl)
	panicIfNotNil(err)

	err = txn.Commit()
	panicIfNotNil(err)

	return &corepb.CreateShortUrlResponse{
		ShortUrl: shortUrl,
	}, nil
}

func (c *ShortUrlsCore) getShortUrl(txn *monstera.Txn, shortUrlId *corepb.ShortUrlId) (*corepb.ShortUrl, error) {
	return c.shortUrlsTable.Get(txn, shortUrlsTablePK(shortUrlId.UserId), shortUrlsTableSK(shortUrlId))
}

func (c *ShortUrlsCore) listShortUrls(txn *monstera.Txn, userId uint32) ([]*corepb.ShortUrl, error) {
	result := make([]*corepb.ShortUrl, 0)

	err := c.shortUrlsTable.List(txn, shortUrlsTablePK(userId), func(shortUrl *corepb.ShortUrl) (bool, error) {
		result = append(result, shortUrl)
		return true, nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *ShortUrlsCore) updateShortUrl(txn *monstera.Txn, shortUrl *corepb.ShortUrl) error {
	return c.shortUrlsTable.Set(txn, shortUrlsTablePK(shortUrl.Id.UserId), shortUrlsTableSK(shortUrl.Id), shortUrl)
}

func (c *ShortUrlsCore) deleteShortUrl(txn *monstera.Txn, shortUrlId *corepb.ShortUrlId) error {
	return c.shortUrlsTable.Delete(txn, shortUrlsTablePK(shortUrlId.UserId), shortUrlsTableSK(shortUrlId))
}

func (c *ShortUrlsCore) createShortUrl(txn *monstera.Txn, shortUrl *corepb.ShortUrl) error {
	return c.shortUrlsTable.Set(txn, shortUrlsTablePK(shortUrl.Id.UserId), shortUrlsTableSK(shortUrl.Id), shortUrl)
}

// 1. shard key (by user id)
// 2. user id
func shortUrlsTablePK(userId uint32) []byte {
	return monstera.ConcatBytes(shardByUser(userId), userId)
}

// 1. short url id
func shortUrlsTableSK(shortUrlId *corepb.ShortUrlId) []byte {
	return monstera.ConcatBytes(shortUrlId.ShortUrlId)
}
