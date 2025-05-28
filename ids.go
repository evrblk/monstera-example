package monsteraexample

import (
	"encoding/binary"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var (
	ErrInvalidId = errors.New("invalid id")

	accountIdRegex = regexp.MustCompile("[0-9a-f]{16}")
)

func DecodeAccountId(s string) (uint64, error) {
	if !accountIdRegex.MatchString(s) {
		return 0, ErrInvalidId
	}

	id, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		return 0, ErrInvalidId
	}

	return id, nil
}

func EncodeAccountId(id uint64) string {
	src := make([]byte, 8)
	binary.BigEndian.PutUint64(src[0:], id)
	return fmt.Sprintf("%x", src)
}
