package common

import "sync"

var mutex = &sync.Mutex{}

// Lock locks the mutex
func Lock() {
	mutex.Lock()
}

// Unlock unlocks the mutex
func UnLock() {
	mutex.Unlock()
}
