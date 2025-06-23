package dlocks

import (
	"io"

	"errors"

	"github.com/evrblk/monstera"
	"github.com/evrblk/monstera-example/dlocks/corepb"
	monsterax "github.com/evrblk/monstera/x"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
)

type LocksCore struct {
	badgerStore *monstera.BadgerStore
	locksTable  *monsterax.CompositeKeyTable[*corepb.Lock, corepb.Lock]
}

var _ LocksCoreApi = &LocksCore{}

func NewLocksCore(badgerStore *monstera.BadgerStore, shardLowerBound []byte, shardUpperBound []byte) *LocksCore {
	return &LocksCore{
		badgerStore: badgerStore,
		locksTable:  monsterax.NewCompositeKeyTable[*corepb.Lock, corepb.Lock](locksTableId, shardLowerBound, shardUpperBound),
	}
}

func (c *LocksCore) ranges() []monstera.KeyRange {
	return []monstera.KeyRange{
		c.locksTable.GetTableKeyRange(),
	}
}

func (c *LocksCore) Snapshot() monstera.ApplicationCoreSnapshot {
	return monsterax.Snapshot(c.badgerStore, c.ranges())
}

func (c *LocksCore) Restore(reader io.ReadCloser) error {
	return monsterax.Restore(c.badgerStore, c.ranges(), reader)
}

func (c *LocksCore) Close() {

}

func (c *LocksCore) GetLock(request *corepb.GetLockRequest) (*corepb.GetLockResponse, error) {
	txn := c.badgerStore.Update()
	defer txn.Discard()

	lock, err := c.getLock(txn, request.LockId)
	if err != nil {
		if errors.Is(err, monstera.ErrNotFound) {
			// No lock exists, return an unlocked lock
			return &corepb.GetLockResponse{
				Lock: &corepb.Lock{
					Id:       request.LockId,
					State:    corepb.LockState_UNLOCKED,
					LockedAt: 0,
				},
			}, nil
		} else {
			panic(err)
		}
	}

	// Check expiration
	lock = c.checkLockExpiration(lock, request.Now)
	if lock.State == corepb.LockState_UNLOCKED {
		// Lock is expired, delete it
		err = c.deleteLock(txn, lock)
		panicIfNotNil(err)
	} else {
		// Lock is still held, update unexpired holders
		err = c.updateLock(txn, lock)
		panicIfNotNil(err)
	}

	err = txn.Commit()
	panicIfNotNil(err)

	return &corepb.GetLockResponse{
		Lock: c.checkLockExpiration(lock, request.Now),
	}, nil
}

func (c *LocksCore) DeleteLock(request *corepb.DeleteLockRequest) (*corepb.DeleteLockResponse, error) {
	txn := c.badgerStore.Update()
	defer txn.Discard()

	lock, err := c.getLock(txn, request.LockId)
	if err != nil {
		if errors.Is(err, monstera.ErrNotFound) {
			// No lock exists, do nothing
			return &corepb.DeleteLockResponse{}, nil
		} else {
			panic(err)
		}
	}

	err = c.deleteLock(txn, lock)
	panicIfNotNil(err)

	err = txn.Commit()
	panicIfNotNil(err)

	return &corepb.DeleteLockResponse{}, nil
}

func (c *LocksCore) AcquireLock(request *corepb.AcquireLockRequest) (*corepb.AcquireLockResponse, error) {
	txn := c.badgerStore.Update()
	defer txn.Discard()

	// TODO check total number of locks

	lock, err := c.getLock(txn, request.LockId)
	if err != nil {
		if errors.Is(err, monstera.ErrNotFound) {
			// No lock exists, create a new one
			lock = &corepb.Lock{
				Id:       request.LockId,
				State:    corepb.LockState_UNLOCKED,
				LockedAt: 0,
			}
		} else {
			panic(err)
		}
	} else {
		// Lock exists, lets check if it expired
		lock = c.checkLockExpiration(lock, request.Now)
	}

	lockHolder := &corepb.LockHolder{
		ProcessId: request.ProcessId,
		LockedAt:  request.Now,
		ExpiresAt: request.ExpiresAt,
	}

	switch lock.State {
	case corepb.LockState_UNLOCKED:
		if request.WriteLock {
			// Lock for writes
			lock.State = corepb.LockState_WRITE_LOCKED
			lock.WriteLockHolder = lockHolder
		} else {
			// Lock for reads only
			lock.State = corepb.LockState_READ_LOCKED
			lock.ReadLockHolders = []*corepb.LockHolder{lockHolder}
		}
		lock.LockedAt = request.Now
	case corepb.LockState_READ_LOCKED:
		if request.WriteLock {
			return &corepb.AcquireLockResponse{
				Lock:    lock,
				Success: false, // Already locked for reads, cannot be locked for writes.
			}, nil
		} else {
			// Already locked for reads.
			// Check if the same process_id already holds the lock here.
			existingHolder, ok := lo.Find(lock.ReadLockHolders, func(h *corepb.LockHolder) bool {
				return h.ProcessId == request.ProcessId
			})
			if ok {
				// Update expiration time (extend lock)
				existingHolder.ExpiresAt = request.ExpiresAt
				existingHolder.LockedAt = request.Now
			} else {
				// Add the new lock holder
				lock.ReadLockHolders = append(lock.ReadLockHolders, lockHolder)
			}
		}
	case corepb.LockState_WRITE_LOCKED:
		if request.WriteLock {
			// Already locked for writes. Check if the same process_id already holds the lock here.
			if lock.WriteLockHolder.ProcessId == request.ProcessId {
				// This process_id already holds the lock, repeated locks are considered successful
				// Update expiration time (extend lock)
				lock.WriteLockHolder.ExpiresAt = request.ExpiresAt
				lock.WriteLockHolder.LockedAt = request.Now
			} else {
				return &corepb.AcquireLockResponse{
					Lock:    lock,
					Success: false, // The lock is held by another process
				}, nil
			}

		} else {
			return &corepb.AcquireLockResponse{
				Lock:    lock,
				Success: false, // Already locked for writes, cannot be locked for reads.
			}, nil

		}
	default:
		panic("invalid lock state")
	}

	err = c.updateLock(txn, lock)
	panicIfNotNil(err)

	err = txn.Commit()
	panicIfNotNil(err)

	return &corepb.AcquireLockResponse{
		Lock:    lock,
		Success: true, // Locked successfully by the given process_id
	}, nil
}

func (c *LocksCore) ReleaseLock(request *corepb.ReleaseLockRequest) (*corepb.ReleaseLockResponse, error) {
	txn := c.badgerStore.Update()
	defer txn.Discard()

	lock, err := c.getLock(txn, request.LockId)
	if err != nil {
		if errors.Is(err, monstera.ErrNotFound) {
			// No lock exists, return an unlocked lock
			return &corepb.ReleaseLockResponse{
				Lock: &corepb.Lock{
					Id:       request.LockId,
					State:    corepb.LockState_UNLOCKED,
					LockedAt: 0,
				},
			}, nil
		} else {
			panic(err)
		}
	} else {
		// Lock exists, lets check if it expired
		lock = c.checkLockExpiration(lock, request.Now)
	}

	switch lock.State {
	case corepb.LockState_UNLOCKED:
		// Lock has expired, delete it
		err = c.deleteLock(txn, lock)
		panicIfNotNil(err)
	case corepb.LockState_READ_LOCKED:
		// Remove the holder
		lock.ReadLockHolders = lo.Filter(lock.ReadLockHolders, func(h *corepb.LockHolder, _ int) bool {
			return h.ProcessId != request.ProcessId
		})

		// If no read lock holders left
		if len(lock.ReadLockHolders) == 0 {
			// Unlock
			lock.LockedAt = 0
			lock.State = corepb.LockState_UNLOCKED
			lock.ReadLockHolders = nil

			// Delete it
			err = c.deleteLock(txn, lock)
			panicIfNotNil(err)
		}
	case corepb.LockState_WRITE_LOCKED:
		if lock.WriteLockHolder.ProcessId == request.ProcessId {
			// Unlock
			lock.State = corepb.LockState_UNLOCKED
			lock.LockedAt = 0
			lock.WriteLockHolder = nil

			// Delete it
			err = c.deleteLock(txn, lock)
			panicIfNotNil(err)
		}
	default:
		panic("invalid lock state")
	}

	err = c.updateLock(txn, lock)
	panicIfNotNil(err)

	err = txn.Commit()
	panicIfNotNil(err)

	return &corepb.ReleaseLockResponse{
		Lock: lock,
	}, nil
}

// checkLockExpiration ensures that the lock is still held at the moment `now`. Returns an updated copy of the lock.
func (c *LocksCore) checkLockExpiration(lock *corepb.Lock, now int64) *corepb.Lock {
	result := proto.Clone(lock).(*corepb.Lock)

	switch lock.State {
	case corepb.LockState_UNLOCKED:
		// Lock is unlocked, return as is
		return result
	case corepb.LockState_READ_LOCKED:
		result.ReadLockHolders = lo.Filter(result.ReadLockHolders, func(h *corepb.LockHolder, _ int) bool {
			return h.ExpiresAt >= now
		})
		if len(result.ReadLockHolders) == 0 {
			result.State = corepb.LockState_UNLOCKED
			result.ReadLockHolders = nil
			result.LockedAt = 0
		}
	case corepb.LockState_WRITE_LOCKED:
		if result.WriteLockHolder.ExpiresAt <= now {
			result.State = corepb.LockState_UNLOCKED
			result.WriteLockHolder = nil
			result.LockedAt = 0
		}
	default:
		panic("invalid lock state")
	}

	return result
}

func (c *LocksCore) getLock(txn *monstera.Txn, lockId *corepb.LockId) (*corepb.Lock, error) {
	return c.locksTable.Get(txn, locksTablePK(lockId), locksTableSK(lockId))
}

func (c *LocksCore) updateLock(txn *monstera.Txn, lock *corepb.Lock) error {
	return c.locksTable.Set(txn, locksTablePK(lock.Id), locksTableSK(lock.Id), lock)
}

func (c *LocksCore) deleteLock(txn *monstera.Txn, lock *corepb.Lock) error {
	return c.locksTable.Delete(txn, locksTablePK(lock.Id), locksTableSK(lock.Id))
}

type locksIdIntf interface {
	GetAccountId() uint64
	GetNamespaceName() string
	GetLockName() string
}

// 1. shard key (by account id and namespace name)
// 2. account id
// 3. namespace name
func locksTablePK(n namespaceIdIntf) []byte {
	return monstera.ConcatBytes(shardByAccountAndNamespace(n.GetAccountId(), n.GetNamespaceName()), n.GetAccountId(), n.GetNamespaceName())
}

// 1. lock name
func locksTableSK(l locksIdIntf) []byte {
	return monstera.ConcatBytes(l.GetLockName())
}
