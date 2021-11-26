package common

import (
	"context"
	"fmt"
	"github.com/werf/lockgate"
	"github.com/werf/lockgate/pkg/distributed_locker"
	"golang.org/x/sync/semaphore"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"strings"
	"sync"
)

// Locker is an interface that provides functions to lock resources
type Locker interface {
	// Lock locks the specified resource. Will block until the lock is acquired
	Lock(key string) (string, error)

	// TryLock tries to acquire the lock. If it is currently blocked, the function will return an error
	TryLock(key string) (bool, string, error)

	// Unlock unlocks the specified resource
	Unlock(key string) error
}

// SyncMutexLocker locks resources using Golang's sync package
type SyncMutexLocker struct {
	mutex sync.Mutex
	locks map[string]*semaphore.Weighted
}

var syncMutexLockerInstance *SyncMutexLocker
var syncMutexLockerOnce sync.Once

// GetSyncMutexLockerInstance returns the SyncMutexLocker singleton instance
func GetSyncMutexLockerInstance() *SyncMutexLocker {
	syncMutexLockerOnce.Do(func() {
		syncMutexLockerInstance = &SyncMutexLocker{
			mutex: sync.Mutex{},
			locks: map[string]*semaphore.Weighted{},
		}
	})
	return syncMutexLockerInstance
}

func (sml *SyncMutexLocker) Lock(key string) (string, error) {
	sml.ensureLockKeyExists(key)
	err := sml.locks[key].Acquire(context.TODO(), 1)
	if err != nil {
		return "", err
	}
	return key, nil
}

func (sml *SyncMutexLocker) TryLock(key string) (bool, string, error) {
	sml.ensureLockKeyExists(key)
	acquired := sml.locks[key].TryAcquire(1)
	return acquired, key, nil
}

func (sml *SyncMutexLocker) Unlock(key string) error {
	sml.ensureLockKeyExists(key)
	sml.locks[key].Release(1)
	return nil
}

func (sml *SyncMutexLocker) ensureLockKeyExists(key string) {
	if sml.locks[key] == nil {
		sml.mutex.Lock()
		sml.locks[key] = semaphore.NewWeighted(1)
		sml.mutex.Unlock()
	}
}

////

var mutex = &sync.Mutex{}

var k8sDistributedLockerInstance *K8sDistributedLocker
var k8sDistributedLockerOnce sync.Once

type K8sDistributedLocker struct {
	locker lockgate.Locker
}

func GetK8sDistributedLockerInstance(client dynamic.Interface) *K8sDistributedLocker {
	k8sDistributedLockerOnce.Do(func() {
		// Initialize kubeDynamicClient from https://github.com/kubernetes/client-go.
		locker := distributed_locker.NewKubernetesLocker(
			client, schema.GroupVersionResource{
				Group:    "",
				Version:  "v1",
				Resource: "configmaps",
			}, "sc-locks", GetKeptnNamespace(),
		)
		k8sDistributedLockerInstance = &K8sDistributedLocker{locker: locker}
	})
	return k8sDistributedLockerInstance
}

func (kdl *K8sDistributedLocker) Lock(key string) (string, error) {
	_, lock, err := kdl.locker.Acquire(key, lockgate.AcquireOptions{})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%s", lock.LockName, lock.UUID), nil
}

func (kdl *K8sDistributedLocker) TryLock(key string) (bool, string, error) {
	acquired, lock, err := kdl.locker.Acquire(key, lockgate.AcquireOptions{NonBlocking: true})
	if err != nil {
		return false, "", err
	}
	return acquired, fmt.Sprintf("%s:%s", lock.LockName, lock.UUID), nil
}

func (kdl *K8sDistributedLocker) Unlock(key string) error {
	split := strings.Split(key, ":")
	if len(split) != 2 {
		return fmt.Errorf("invalid lock key %s. Expected <name:uuid>", key)
	}
	err := kdl.locker.Release(lockgate.LockHandle{
		LockName: split[0],
		UUID:     split[1],
	})
	if err != nil {
		return fmt.Errorf("could not release lock: %v", err.Error())
	}
	return nil
}

var projectLocks = map[string]*sync.Mutex{}

// Lock locks the mutex
func Lock() {
	mutex.Lock()
}

// Unlock unlocks the mutex
func Unlock() {
	mutex.Unlock()
}

// LockProject
func LockProject(project string) {
	if projectLocks[project] == nil {
		Lock()
		projectLocks[project] = &sync.Mutex{}
		Unlock()
	}
	projectLocks[project].Lock()
}

func UnlockProject(project string) {
	if projectLocks[project] == nil {
		Lock()
		projectLocks[project] = &sync.Mutex{}
		Unlock()
	}
	projectLocks[project].Unlock()
}

func LockServiceInStageOfProject(project, stage, service string) {
	lockKey := fmt.Sprintf("%s.%s.%s", project, stage, service)
	if projectLocks[lockKey] == nil {
		Lock()
		projectLocks[lockKey] = &sync.Mutex{}
		Unlock()
	}
	projectLocks[lockKey].Lock()
}

func UnlockServiceInStageOfProject(project, stage, service string) {
	lockKey := fmt.Sprintf("%s.%s.%s", project, stage, service)
	if projectLocks[lockKey] == nil {
		Lock()
		projectLocks[lockKey] = &sync.Mutex{}
		Unlock()
	}
	projectLocks[lockKey].Unlock()
}
