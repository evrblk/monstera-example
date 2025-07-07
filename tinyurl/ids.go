package tinyurl

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/evrblk/monstera-example/tinyurl/base62"
	"github.com/evrblk/monstera-example/tinyurl/corepb"
	"regexp"
)

var (
	ErrInvalidId = errors.New("invalid id")

	shortUrlIdRegex = regexp.MustCompile("[0-9a-zA-Z]+")
	userIdRegex     = regexp.MustCompile("[0-9a-zA-Z]+")
)

func DecodeUserId(s string) (uint32, error) {
	if !userIdRegex.MatchString(s) {
		return 0, ErrInvalidId
	}

	b, err := base62.DecodeString(s)
	if err != nil {
		return 0, ErrInvalidId
	}

	if len(b) != 4 {
		return 0, ErrInvalidId
	}

	return binary.BigEndian.Uint32(b[0:4]), nil
}

func EncodeUserId(id uint32) string {
	src := make([]byte, 4)
	binary.BigEndian.PutUint32(src[0:], id)
	return fmt.Sprintf("%s", base62.Encode(src))
}

func DecodeShortUrlId(s string) (*corepb.ShortUrlId, error) {
	if !shortUrlIdRegex.MatchString(s) {
		return nil, ErrInvalidId
	}

	b, err := base62.DecodeString(s)
	if err != nil {
		return nil, ErrInvalidId
	}

	if len(b) != 4+4 {
		return nil, ErrInvalidId
	}

	return &corepb.ShortUrlId{
		UserId:     binary.BigEndian.Uint32(b[0:4]),
		ShortUrlId: binary.BigEndian.Uint32(b[4:8]),
	}, nil
}

func EncodeShortUrlId(id *corepb.ShortUrlId) string {
	src := make([]byte, 4+4)
	binary.BigEndian.PutUint32(src[0:4], id.UserId)
	binary.BigEndian.PutUint32(src[4:8], id.ShortUrlId)
	return fmt.Sprintf("%s", base62.Encode(src))
}
