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
	Lock()
	if projectLocks[project] == nil {
		projectLocks[project] = &sync.Mutex{}
	}
	Unlock()
	projectLocks[project].Lock()
}

func UnlockProject(project string) {
	Lock()
	if projectLocks[project] == nil {
		projectLocks[project] = &sync.Mutex{}
	}
	Unlock()
	projectLocks[project].Unlock()
}

func LockServiceInStageOfProject(project, stage, service string) {
	lockKey := fmt.Sprintf("%s.%s.%s", project, stage, service)
	Lock()
	if projectLocks[lockKey] == nil {
		projectLocks[lockKey] = &sync.Mutex{}
	}
	Unlock()
	projectLocks[lockKey].Lock()
}

func UnlockServiceInStageOfProject(project, stage, service string) {
	lockKey := fmt.Sprintf("%s.%s.%s", project, stage, service)
	Lock()
	if projectLocks[lockKey] == nil {
		projectLocks[lockKey] = &sync.Mutex{}
	}
	Unlock()
	projectLocks[lockKey].Unlock()
}
