package tinyurl

import (
	"github.com/evrblk/monstera-example/tinyurl/corepb"
	"github.com/evrblk/monstera-example/tinyurl/gatewaypb"
)

func userToFront(user *corepb.User) *gatewaypb.User {
	if user == nil {
		return nil
	}

	return &gatewaypb.User{
		Id:        EncodeUserId(user.Id),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func shortUrlToFront(shortUrl *corepb.ShortUrl) *gatewaypb.ShortUrl {
	if shortUrl == nil {
		return nil
	}

	return &gatewaypb.ShortUrl{
		Id:        EncodeShortUrlId(shortUrl.Id),
		FullUrl:   shortUrl.FullUrl,
		CreatedAt: shortUrl.CreatedAt,
	}
}
