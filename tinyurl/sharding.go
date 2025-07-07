package tinyurl

import (
	"github.com/evrblk/monstera"
	"github.com/evrblk/monstera-example/tinyurl/corepb"
)

func shardByUser(userId uint32) []byte {
	return monstera.GetShardKey(monstera.ConcatBytes(userId), 4)
}

type ShardKeyCalculator struct{}

var _ TinyUrlServiceMonsteraShardKeyCalculator = &ShardKeyCalculator{}

func (s ShardKeyCalculator) GetShortUrlShardKey(request *corepb.GetShortUrlRequest) []byte {
	return shardByUser(request.ShortUrlId.UserId)
}

func (s ShardKeyCalculator) ListShortUrlsShardKey(request *corepb.ListShortUrlsRequest) []byte {
	return shardByUser(request.UserId)
}

func (s ShardKeyCalculator) CreateShortUrlShardKey(request *corepb.CreateShortUrlRequest) []byte {
	return shardByUser(request.ShortUrlId.UserId)
}

func (s ShardKeyCalculator) GetUserShardKey(request *corepb.GetUserRequest) []byte {
	return shardByUser(request.UserId)
}

func (s ShardKeyCalculator) CreateUserShardKey(request *corepb.CreateUserRequest) []byte {
	return shardByUser(request.UserId)
}
