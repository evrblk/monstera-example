package monsteraexample

import (
	"math/rand/v2"
	"testing"
	"time"

	"github.com/evrblk/monstera"
	"github.com/evrblk/monstera-example/corepb"
	"github.com/stretchr/testify/require"
)

func TestAcquireWriteLock(t *testing.T) {
	require := require.New(t)

	locksCore := newLocksCore()

	now := time.Now()

	accountId := rand.Uint64()
	lockId := &corepb.LockId{
		AccountId:     accountId,
		NamespaceName: "test_namespace",
		LockName:      "test_lock",
	}

	// T+0: Acquire lock
	response1, err := locksCore.AcquireLock(&corepb.AcquireLockRequest{
		LockId:    lockId,
		Now:       now.UnixNano(),
		ProcessId: "process_1",
		WriteLock: true,
		ExpiresAt: now.Add(time.Hour).UnixNano(),
	})

	require.NoError(err)
	require.NotNil(response1.Lock)
	require.Equal(true, response1.Success)
	require.EqualValues(now.UnixNano(), response1.Lock.LockedAt)
	require.Equal(corepb.LockState_WRITE_LOCKED, response1.Lock.State)
	require.Equal("process_1", response1.Lock.WriteLockHolder.ProcessId)
	require.EqualValues(now.Add(time.Hour).UnixNano(), response1.Lock.WriteLockHolder.ExpiresAt)
	require.EqualValues(now.UnixNano(), response1.Lock.WriteLockHolder.LockedAt)

	// T+1m: Get lock
	response2, err := locksCore.GetLock(&corepb.GetLockRequest{
		LockId: lockId,
		Now:    now.Add(time.Minute).UnixNano(),
	})

	require.NoError(err)
	require.NotNil(response2.Lock)
	require.Equal(corepb.LockState_WRITE_LOCKED, response1.Lock.State)

	// T+61m: Get lock
	response3, err := locksCore.GetLock(&corepb.GetLockRequest{
		LockId: lockId,
		Now:    now.Add(61 * time.Minute).UnixNano(),
	})

	require.NoError(err)
	require.NotNil(response3.Lock)
	require.Equal(corepb.LockState_UNLOCKED, response3.Lock.State)
}

func TestAcquireReadLock(t *testing.T) {
	require := require.New(t)

	locksCore := newLocksCore()

	now := time.Now()

	accountId := rand.Uint64()
	lockId := &corepb.LockId{
		AccountId:     accountId,
		NamespaceName: "test_namespace",
		LockName:      "test_lock",
	}

	// T+0: Acquire lock
	response1, err := locksCore.AcquireLock(&corepb.AcquireLockRequest{
		LockId:    lockId,
		Now:       now.UnixNano(),
		ProcessId: "process_1",
		WriteLock: false,
		ExpiresAt: now.Add(time.Hour).UnixNano(),
	})

	require.NoError(err)
	require.NotNil(response1.Lock)
	require.Equal(true, response1.Success)
	require.EqualValues(now.UnixNano(), response1.Lock.LockedAt)
	require.Equal(corepb.LockState_READ_LOCKED, response1.Lock.State)
	require.Len(response1.Lock.ReadLockHolders, 1)
	require.Equal("process_1", response1.Lock.ReadLockHolders[0].ProcessId)
	require.EqualValues(now.Add(time.Hour).UnixNano(), response1.Lock.ReadLockHolders[0].ExpiresAt)
	require.EqualValues(now.UnixNano(), response1.Lock.ReadLockHolders[0].LockedAt)

	// T+1m: Get lock
	response2, err := locksCore.GetLock(&corepb.GetLockRequest{
		LockId: lockId,
		Now:    now.Add(time.Minute).UnixNano(),
	})

	require.NoError(err)
	require.NotNil(response2.Lock)
	require.Equal(corepb.LockState_READ_LOCKED, response1.Lock.State)

	// T+61m: Get lock
	response3, err := locksCore.GetLock(&corepb.GetLockRequest{
		LockId: lockId,
		Now:    now.Add(61 * time.Minute).UnixNano(),
	})

	require.NoError(err)
	require.NotNil(response3.Lock)
	require.Equal(corepb.LockState_UNLOCKED, response3.Lock.State)
}

func TestAcquireWriteLockRepeatedly(t *testing.T) {
	require := require.New(t)

	locksCore := newLocksCore()

	now := time.Now()

	accountId := rand.Uint64()
	lockId := &corepb.LockId{
		AccountId:     accountId,
		NamespaceName: "test_namespace",
		LockName:      "test_lock",
	}

	// T+0: Acquire lock
	response1, err := locksCore.AcquireLock(&corepb.AcquireLockRequest{
		LockId:    lockId,
		Now:       now.UnixNano(),
		ProcessId: "process_1",
		WriteLock: true,
		ExpiresAt: now.Add(time.Hour).UnixNano(),
	})

	require.NoError(err)
	require.Equal(true, response1.Success)
	require.Equal(corepb.LockState_WRITE_LOCKED, response1.Lock.State)

	// T+1m: Acquire lock again
	response2, err := locksCore.AcquireLock(&corepb.AcquireLockRequest{
		LockId:    lockId,
		Now:       now.Add(time.Minute).UnixNano(),
		ProcessId: "process_1",
		WriteLock: true,
		ExpiresAt: now.Add(time.Minute).Add(time.Hour).UnixNano(),
	})

	require.NoError(err)
	require.Equal(true, response2.Success)
	require.Equal(corepb.LockState_WRITE_LOCKED, response2.Lock.State)
	require.EqualValues(now.Add(time.Minute).Add(time.Hour).UnixNano(), response2.Lock.WriteLockHolder.ExpiresAt)
}

func TestAcquireReadLockRepeatedly(t *testing.T) {
	require := require.New(t)

	locksCore := newLocksCore()

	now := time.Now()

	accountId := rand.Uint64()
	lockId := &corepb.LockId{
		AccountId:     accountId,
		NamespaceName: "test_namespace",
		LockName:      "test_lock",
	}

	// T+0: Acquire lock
	response1, err := locksCore.AcquireLock(&corepb.AcquireLockRequest{
		LockId:    lockId,
		Now:       now.UnixNano(),
		ProcessId: "process_1",
		WriteLock: false,
		ExpiresAt: now.Add(time.Hour).UnixNano(),
	})

	require.NoError(err)
	require.Equal(true, response1.Success)
	require.Equal(corepb.LockState_READ_LOCKED, response1.Lock.State)

	// T+1m: Acquire lock again
	response2, err := locksCore.AcquireLock(&corepb.AcquireLockRequest{
		LockId:    lockId,
		Now:       now.Add(time.Minute).UnixNano(),
		ProcessId: "process_1",
		WriteLock: false,
		ExpiresAt: now.Add(time.Minute).Add(time.Hour).UnixNano(),
	})

	require.NoError(err)
	require.Equal(true, response2.Success)
	require.Equal(corepb.LockState_READ_LOCKED, response2.Lock.State)
	require.EqualValues(now.Add(time.Minute).Add(time.Hour).UnixNano(), response2.Lock.ReadLockHolders[0].ExpiresAt)
}

func TestAcquireLockWriteLockedByAnotherProcess(t *testing.T) {
	require := require.New(t)

	locksCore := newLocksCore()

	now := time.Now()

	accountId := rand.Uint64()
	lockId := &corepb.LockId{
		AccountId:     accountId,
		NamespaceName: "test_namespace",
		LockName:      "test_lock",
	}

	// T+0: Acquire lock
	response1, err := locksCore.AcquireLock(&corepb.AcquireLockRequest{
		LockId:    lockId,
		Now:       now.UnixNano(),
		ProcessId: "process_1",
		WriteLock: true,
		ExpiresAt: now.Add(time.Hour).UnixNano(),
	})

	require.NoError(err)
	require.Equal(true, response1.Success)
	require.Equal(corepb.LockState_WRITE_LOCKED, response1.Lock.State)

	// T+1m: Acquire write lock by another process
	response2, err := locksCore.AcquireLock(&corepb.AcquireLockRequest{
		LockId:    lockId,
		Now:       now.Add(time.Minute).UnixNano(),
		ProcessId: "process_2",
		WriteLock: true,
		ExpiresAt: now.Add(time.Minute).Add(time.Hour).UnixNano(),
	})

	require.NoError(err)
	require.Equal(false, response2.Success)
	require.Equal(corepb.LockState_WRITE_LOCKED, response2.Lock.State)
	require.Equal("process_1", response2.Lock.WriteLockHolder.ProcessId)

	// T+2m: Acquire read lock by another process
	response3, err := locksCore.AcquireLock(&corepb.AcquireLockRequest{
		LockId:    lockId,
		Now:       now.Add(2 * time.Minute).UnixNano(),
		ProcessId: "process_2",
		WriteLock: false,
		ExpiresAt: now.Add(2 * time.Minute).Add(time.Hour).UnixNano(),
	})

	require.NoError(err)
	require.Equal(false, response3.Success)
	require.Equal(corepb.LockState_WRITE_LOCKED, response3.Lock.State)
	require.Equal("process_1", response3.Lock.WriteLockHolder.ProcessId)
}

func TestAcquireLockReadLockedByAnotherProcess(t *testing.T) {
	require := require.New(t)

	locksCore := newLocksCore()

	now := time.Now()

	accountId := rand.Uint64()
	lockId := &corepb.LockId{
		AccountId:     accountId,
		NamespaceName: "test_namespace",
		LockName:      "test_lock",
	}

	// T+0: Acquire lock
	response1, err := locksCore.AcquireLock(&corepb.AcquireLockRequest{
		LockId:    lockId,
		Now:       now.UnixNano(),
		ProcessId: "process_1",
		WriteLock: false,
		ExpiresAt: now.Add(time.Hour).UnixNano(),
	})

	require.NoError(err)
	require.Equal(true, response1.Success)
	require.Equal(corepb.LockState_READ_LOCKED, response1.Lock.State)

	// T+1m: Acquire write lock by another process
	response2, err := locksCore.AcquireLock(&corepb.AcquireLockRequest{
		LockId:    lockId,
		Now:       now.Add(time.Minute).UnixNano(),
		ProcessId: "process_2",
		WriteLock: true,
		ExpiresAt: now.Add(time.Minute).Add(time.Hour).UnixNano(),
	})

	require.NoError(err)
	require.Equal(false, response2.Success)
	require.Equal(corepb.LockState_READ_LOCKED, response2.Lock.State)
	require.Equal("process_1", response2.Lock.ReadLockHolders[0].ProcessId)

	// T+2m: Acquire read lock by another process
	response3, err := locksCore.AcquireLock(&corepb.AcquireLockRequest{
		LockId:    lockId,
		Now:       now.Add(2 * time.Minute).UnixNano(),
		ProcessId: "process_2",
		WriteLock: false,
		ExpiresAt: now.Add(2 * time.Minute).Add(time.Hour).UnixNano(),
	})

	require.NoError(err)
	require.Equal(true, response3.Success)
	require.Equal(corepb.LockState_READ_LOCKED, response3.Lock.State)
	require.Len(response3.Lock.ReadLockHolders, 2)
	require.Equal("process_1", response3.Lock.ReadLockHolders[0].ProcessId)
	require.Equal("process_2", response3.Lock.ReadLockHolders[1].ProcessId)
}

func TestGetNonexistentLock(t *testing.T) {
	require := require.New(t)

	locksCore := newLocksCore()

	now := time.Now()

	accountId := rand.Uint64()
	lockId := &corepb.LockId{
		AccountId:     accountId,
		NamespaceName: "test_namespace",
		LockName:      "test_lock",
	}

	// Get lock
	response1, err := locksCore.GetLock(&corepb.GetLockRequest{
		LockId: lockId,
		Now:    now.UnixNano(),
	})

	require.NoError(err)
	require.NotNil(response1.Lock)
	require.Equal(corepb.LockState_UNLOCKED, response1.Lock.State)
}

func TestDeleteNonexistentLock(t *testing.T) {
	require := require.New(t)

	locksCore := newLocksCore()

	now := time.Now()

	accountId := rand.Uint64()
	lockId := &corepb.LockId{
		AccountId:     accountId,
		NamespaceName: "test_namespace",
		LockName:      "test_lock",
	}

	// Delete lock
	_, err := locksCore.DeleteLock(&corepb.DeleteLockRequest{
		LockId: lockId,
		Now:    now.UnixNano(),
	})

	require.NoError(err)
}

func TestReleaseNonexistentLock(t *testing.T) {
	require := require.New(t)

	locksCore := newLocksCore()

	now := time.Now()

	accountId := rand.Uint64()
	lockId := &corepb.LockId{
		AccountId:     accountId,
		NamespaceName: "test_namespace",
		LockName:      "test_lock",
	}

	// Release lock
	response1, err := locksCore.ReleaseLock(&corepb.ReleaseLockRequest{
		LockId:    lockId,
		ProcessId: "process_1",
		Now:       now.UnixNano(),
	})

	require.NoError(err)
	require.NotNil(response1.Lock)
	require.Equal(corepb.LockState_UNLOCKED, response1.Lock.State)
}

func TestDeleteAcquiredLock(t *testing.T) {
	require := require.New(t)

	locksCore := newLocksCore()

	now := time.Now()

	accountId := rand.Uint64()
	lockId := &corepb.LockId{
		AccountId:     accountId,
		NamespaceName: "test_namespace",
		LockName:      "test_lock",
	}

	// T+0: Acquire lock
	response1, err := locksCore.AcquireLock(&corepb.AcquireLockRequest{
		LockId:    lockId,
		Now:       now.UnixNano(),
		ProcessId: "process_1",
		WriteLock: true,
		ExpiresAt: now.Add(time.Hour).UnixNano(),
	})

	require.NoError(err)
	require.NotNil(response1.Lock)
	require.Equal(true, response1.Success)
	require.Equal(corepb.LockState_WRITE_LOCKED, response1.Lock.State)

	// T+1m: Delete lock
	_, err = locksCore.DeleteLock(&corepb.DeleteLockRequest{
		LockId: lockId,
		Now:    now.Add(time.Minute).UnixNano(),
	})

	require.NoError(err)

	// T+2m: Acquire lock
	response3, err := locksCore.AcquireLock(&corepb.AcquireLockRequest{
		LockId:    lockId,
		Now:       now.Add(2 * time.Minute).UnixNano(),
		ProcessId: "process_2",
		WriteLock: true,
		ExpiresAt: now.Add(2 * time.Minute).Add(time.Hour).UnixNano(),
	})

	require.NoError(err)
	require.NotNil(response3.Lock)
	require.Equal(true, response3.Success)
	require.Equal(corepb.LockState_WRITE_LOCKED, response3.Lock.State)
	require.Equal("process_2", response3.Lock.WriteLockHolder.ProcessId)
}

func TestReleaseWriteLock(t *testing.T) {
	require := require.New(t)

	locksCore := newLocksCore()

	now := time.Now()

	accountId := rand.Uint64()
	lockId := &corepb.LockId{
		AccountId:     accountId,
		NamespaceName: "test_namespace",
		LockName:      "test_lock",
	}

	// T+0: Acquire lock
	response1, err := locksCore.AcquireLock(&corepb.AcquireLockRequest{
		LockId:    lockId,
		Now:       now.UnixNano(),
		ProcessId: "process_1",
		WriteLock: true,
		ExpiresAt: now.Add(time.Hour).UnixNano(),
	})

	require.NoError(err)
	require.NotNil(response1.Lock)
	require.Equal(true, response1.Success)
	require.Equal(corepb.LockState_WRITE_LOCKED, response1.Lock.State)

	// T+1m: Release lock with wrong process id
	response2, err := locksCore.ReleaseLock(&corepb.ReleaseLockRequest{
		LockId:    lockId,
		ProcessId: "process_2",
		Now:       now.Add(time.Minute).UnixNano(),
	})

	require.NoError(err)
	require.NotNil(response2.Lock)
	require.Equal(corepb.LockState_WRITE_LOCKED, response1.Lock.State)
	require.Equal("process_1", response1.Lock.WriteLockHolder.ProcessId)

	// T+2m: Release lock with correct process id
	response3, err := locksCore.ReleaseLock(&corepb.ReleaseLockRequest{
		LockId:    lockId,
		Now:       now.Add(2 * time.Minute).UnixNano(),
		ProcessId: "process_1",
	})

	require.NoError(err)
	require.NotNil(response3.Lock)
	require.Equal(corepb.LockState_UNLOCKED, response3.Lock.State)
}

func TestReleaseReadLock(t *testing.T) {
	require := require.New(t)

	locksCore := newLocksCore()

	now := time.Now()

	accountId := rand.Uint64()
	lockId := &corepb.LockId{
		AccountId:     accountId,
		NamespaceName: "test_namespace",
		LockName:      "test_lock",
	}

	// T+0: Acquire read lock
	response1, err := locksCore.AcquireLock(&corepb.AcquireLockRequest{
		LockId:    lockId,
		Now:       now.UnixNano(),
		ProcessId: "process_1",
		WriteLock: false,
		ExpiresAt: now.Add(time.Hour).UnixNano(),
	})

	require.NoError(err)
	require.NotNil(response1.Lock)
	require.Equal(true, response1.Success)
	require.Equal(corepb.LockState_READ_LOCKED, response1.Lock.State)

	// T+1m: Acquire read lock from another process
	response2, err := locksCore.AcquireLock(&corepb.AcquireLockRequest{
		LockId:    lockId,
		Now:       now.Add(time.Minute).UnixNano(),
		ProcessId: "process_2",
		WriteLock: false,
		ExpiresAt: now.Add(time.Minute).Add(time.Hour).UnixNano(),
	})

	require.NoError(err)
	require.NotNil(response2.Lock)
	require.Equal(true, response2.Success)
	require.Equal(corepb.LockState_READ_LOCKED, response2.Lock.State)

	// T+2m: Release lock with first process id
	response3, err := locksCore.ReleaseLock(&corepb.ReleaseLockRequest{
		LockId:    lockId,
		Now:       now.Add(2 * time.Minute).UnixNano(),
		ProcessId: "process_1",
	})

	require.NoError(err)
	require.NotNil(response3.Lock)
	require.Equal(corepb.LockState_READ_LOCKED, response3.Lock.State)
	require.Equal("process_2", response3.Lock.ReadLockHolders[0].ProcessId)

	// T+3m: Release lock with second process id
	response4, err := locksCore.ReleaseLock(&corepb.ReleaseLockRequest{
		LockId:    lockId,
		Now:       now.Add(3 * time.Minute).UnixNano(),
		ProcessId: "process_2",
	})

	require.NoError(err)
	require.NotNil(response4.Lock)
	require.Equal(corepb.LockState_UNLOCKED, response4.Lock.State)
	require.Len(response4.Lock.ReadLockHolders, 0)
}

func TestReleaseExpiredLock(t *testing.T) {
	require := require.New(t)

	locksCore := newLocksCore()

	now := time.Now()

	accountId := rand.Uint64()
	lockId := &corepb.LockId{
		AccountId:     accountId,
		NamespaceName: "test_namespace",
		LockName:      "test_lock",
	}

	// T+0: Acquire lock
	response1, err := locksCore.AcquireLock(&corepb.AcquireLockRequest{
		LockId:    lockId,
		Now:       now.UnixNano(),
		ProcessId: "process_1",
		WriteLock: true,
		ExpiresAt: now.Add(time.Hour).UnixNano(),
	})

	require.NoError(err)
	require.NotNil(response1.Lock)
	require.Equal(true, response1.Success)
	require.Equal(corepb.LockState_WRITE_LOCKED, response1.Lock.State)

	// T+61m: Release lock after expiration time
	response2, err := locksCore.ReleaseLock(&corepb.ReleaseLockRequest{
		LockId:    lockId,
		Now:       now.Add(61 * time.Minute).UnixNano(),
		ProcessId: "process_1",
	})

	require.NoError(err)
	require.NotNil(response2.Lock)
	require.Equal(corepb.LockState_UNLOCKED, response2.Lock.State)
}

func newLocksCore() *LocksCore {
	return NewLocksCore(monstera.NewBadgerInMemoryStore(), []byte{0x00, 0x00}, []byte{0xff, 0xff})
}
