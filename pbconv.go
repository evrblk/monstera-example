package monsteraexample

import (
	"github.com/evrblk/monstera-example/corepb"
	"github.com/evrblk/monstera-example/gatewaypb"
)

func namespaceToFront(namespace *corepb.Namespace) *gatewaypb.Namespace {
	if namespace == nil {
		return nil
	}

	return &gatewaypb.Namespace{
		Name:        namespace.Id.NamespaceName,
		Description: namespace.Description,
		CreatedAt:   namespace.CreatedAt,
		UpdatedAt:   namespace.UpdatedAt,
	}
}

func namespacesToFront(namespaces []*corepb.Namespace) []*gatewaypb.Namespace {
	frontNamespaces := make([]*gatewaypb.Namespace, len(namespaces))
	for i, namespace := range namespaces {
		frontNamespaces[i] = namespaceToFront(namespace)
	}
	return frontNamespaces
}

func lockToFront(lock *corepb.Lock) *gatewaypb.Lock {
	if lock == nil {
		return nil
	}

	return &gatewaypb.Lock{
		Name:            lock.Id.LockName,
		State:           gatewaypb.LockState(lock.State),
		LockedAt:        lock.LockedAt,
		WriteLockHolder: lockHolderToFront(lock.WriteLockHolder),
		ReadLockHolders: lockHoldersToFront(lock.ReadLockHolders),
	}
}

func lockHolderToFront(lockHolder *corepb.LockHolder) *gatewaypb.LockHolder {
	if lockHolder == nil {
		return nil
	}

	return &gatewaypb.LockHolder{
		ProcessId: lockHolder.ProcessId,
		LockedAt:  lockHolder.LockedAt,
		ExpiresAt: lockHolder.ExpiresAt,
	}
}

func lockHoldersToFront(lockHolders []*corepb.LockHolder) []*gatewaypb.LockHolder {
	frontLockHolders := make([]*gatewaypb.LockHolder, len(lockHolders))
	for i, lockHolder := range lockHolders {
		frontLockHolders[i] = lockHolderToFront(lockHolder)
	}
	return frontLockHolders
}
