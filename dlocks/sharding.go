package dlocks

import (
	"github.com/evrblk/monstera"
	"github.com/evrblk/monstera-example/dlocks/corepb"
)

func shardByAccount(accountId uint64) []byte {
	return monstera.GetShardKey(monstera.ConcatBytes(accountId), 4)
}

func shardByAccountAndNamespace(accountId uint64, namespaceName string) []byte {
	return monstera.GetShardKey(monstera.ConcatBytes(accountId, namespaceName), 4)
}

type ShardKeyCalculator struct{}

var _ LocksServiceMonsteraShardKeyCalculator = &ShardKeyCalculator{}

func (g *ShardKeyCalculator) CreateAccountShardKey(request *corepb.CreateAccountRequest) []byte {
	return []byte{0x00, 0x00, 0x00, 0x00}
}

func (g *ShardKeyCalculator) DeleteAccountShardKey(request *corepb.DeleteAccountRequest) []byte {
	return []byte{0x00, 0x00, 0x00, 0x00}
}

func (g *ShardKeyCalculator) GetAccountShardKey(request *corepb.GetAccountRequest) []byte {
	return []byte{0x00, 0x00, 0x00, 0x00}
}

func (g *ShardKeyCalculator) ListAccountsShardKey(request *corepb.ListAccountsRequest) []byte {
	return []byte{0x00, 0x00, 0x00, 0x00}
}

func (g *ShardKeyCalculator) UpdateAccountShardKey(request *corepb.UpdateAccountRequest) []byte {
	return []byte{0x00, 0x00, 0x00, 0x00}
}

func (g *ShardKeyCalculator) AcquireLockShardKey(request *corepb.AcquireLockRequest) []byte {
	return shardByAccountAndNamespace(request.LockId.AccountId, request.LockId.NamespaceName)
}

func (g *ShardKeyCalculator) CreateNamespaceShardKey(request *corepb.CreateNamespaceRequest) []byte {
	return shardByAccount(request.AccountId)
}

func (g *ShardKeyCalculator) DeleteLockShardKey(request *corepb.DeleteLockRequest) []byte {
	return shardByAccountAndNamespace(request.LockId.AccountId, request.LockId.NamespaceName)
}

func (g *ShardKeyCalculator) DeleteNamespaceShardKey(request *corepb.DeleteNamespaceRequest) []byte {
	return shardByAccount(request.NamespaceId.AccountId)
}

func (g *ShardKeyCalculator) GetLockShardKey(request *corepb.GetLockRequest) []byte {
	return shardByAccountAndNamespace(request.LockId.AccountId, request.LockId.NamespaceName)
}

func (g *ShardKeyCalculator) GetNamespaceShardKey(request *corepb.GetNamespaceRequest) []byte {
	return shardByAccount(request.NamespaceId.AccountId)
}

func (g *ShardKeyCalculator) ListNamespacesShardKey(request *corepb.ListNamespacesRequest) []byte {
	return shardByAccount(request.AccountId)
}

func (g *ShardKeyCalculator) ReleaseLockShardKey(request *corepb.ReleaseLockRequest) []byte {
	return shardByAccountAndNamespace(request.LockId.AccountId, request.LockId.NamespaceName)
}

func (g *ShardKeyCalculator) UpdateNamespaceShardKey(request *corepb.UpdateNamespaceRequest) []byte {
	return shardByAccount(request.NamespaceId.AccountId)
}
