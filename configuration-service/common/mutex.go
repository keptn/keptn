package common

import "sync"

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
