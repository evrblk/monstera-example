package ledger

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/evrblk/monstera-example/ledger/corepb"
	"regexp"
	"strconv"
)

var (
	ErrInvalidId = errors.New("invalid id")

	accountIdRegex     = regexp.MustCompile("[0-9a-f]{16}")
	transactionIdRegex = regexp.MustCompile("[0-9a-f]{32}")
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

func DecodeTransactionId(s string) (*corepb.TransactionId, error) {
	if !transactionIdRegex.MatchString(s) {
		return nil, ErrInvalidId
	}

	accountId, err := strconv.ParseUint(s[0:16], 16, 64)
	if err != nil {
		return nil, ErrInvalidId
	}

	transactionId, err := strconv.ParseUint(s[16:32], 16, 64)
	if err != nil {
		return nil, ErrInvalidId
	}

	return &corepb.TransactionId{
		AccountId:     accountId,
		TransactionId: transactionId,
	}, nil
}

func EncodeTransactionId(id *corepb.TransactionId) string {
	src := make([]byte, 8+8)
	binary.BigEndian.PutUint64(src[0:], id.AccountId)
	binary.BigEndian.PutUint64(src[8:], id.TransactionId)
	return fmt.Sprintf("%x", src)
}
