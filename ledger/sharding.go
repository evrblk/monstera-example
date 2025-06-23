package ledger

import (
	"github.com/evrblk/monstera"
	"github.com/evrblk/monstera-example/ledger/corepb"
)

func shardByAccount(accountId uint64) []byte {
	return monstera.GetShardKey(monstera.ConcatBytes(accountId), 4)
}

type ShardKeyCalculator struct{}

var _ LedgerServiceMonsteraShardKeyCalculator = &ShardKeyCalculator{}

func (g *ShardKeyCalculator) CreateAccountShardKey(request *corepb.CreateAccountRequest) []byte {
	return shardByAccount(request.AccountId)
}

func (g *ShardKeyCalculator) ListTransactionsShardKey(request *corepb.ListTransactionsRequest) []byte {
	return shardByAccount(request.AccountId)
}

func (g *ShardKeyCalculator) GetTransactionShardKey(request *corepb.GetTransactionRequest) []byte {
	return shardByAccount(request.TransactionId.AccountId)
}

func (g *ShardKeyCalculator) CreateTransactionShardKey(request *corepb.CreateTransactionRequest) []byte {
	return shardByAccount(request.TransactionId.AccountId)
}

func (g *ShardKeyCalculator) CancelTransactionShardKey(request *corepb.CancelTransactionRequest) []byte {
	return shardByAccount(request.TransactionId.AccountId)
}

func (g *ShardKeyCalculator) SettleTransactionShardKey(request *corepb.SettleTransactionRequest) []byte {
	return shardByAccount(request.TransactionId.AccountId)
}

func (g *ShardKeyCalculator) GetAccountShardKey(request *corepb.GetAccountRequest) []byte {
	return shardByAccount(request.AccountId)
}
