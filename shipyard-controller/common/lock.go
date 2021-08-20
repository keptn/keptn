package common

import (
	"fmt"
	"sync"
)

var mutex = &sync.Mutex{}

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
