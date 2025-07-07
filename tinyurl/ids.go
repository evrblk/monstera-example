package tinyurl

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/evrblk/monstera-example/tinyurl/corepb"
	"regexp"
	"strconv"
)

var (
	ErrInvalidId = errors.New("invalid id")

	userIdRegex = regexp.MustCompile("[0-9a-f]{8}")
)

func DecodeUserId(s string) (uint32, error) {
	if !userIdRegex.MatchString(s) {
		return 0, ErrInvalidId
	}

	id, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return 0, ErrInvalidId
	}

	return uint32(id), nil
}

func EncodeUserId(id uint32) string {
	src := make([]byte, 4)
	binary.BigEndian.PutUint32(src[0:], id)
	return fmt.Sprintf("%x", src)
}

func DecodeShortUrlId(s string) (*corepb.ShortUrlId, error) {
	//if !transactionIdRegex.MatchString(s) {
	//	return nil, ErrInvalidId
	//}

	userId, err := strconv.ParseUint(s[0:8], 16, 32)
	if err != nil {
		return nil, ErrInvalidId
	}

	transactionId, err := strconv.ParseUint(s[8:16], 16, 32)
	if err != nil {
		return nil, ErrInvalidId
	}

	return &corepb.ShortUrlId{
		UserId:     uint32(userId),
		ShortUrlId: uint32(transactionId),
	}, nil
}

func EncodeShortUrlId(id *corepb.ShortUrlId) string {
	src := make([]byte, 4+4)
	binary.BigEndian.PutUint32(src[0:], id.UserId)
	binary.BigEndian.PutUint32(src[4:], id.ShortUrlId)
	return fmt.Sprintf("%x", src)
}
