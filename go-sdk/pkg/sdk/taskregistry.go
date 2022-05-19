package sdk

import (
	"sync"
)

type taskRegistry struct {
	sync.RWMutex
	entries map[string]taskEntry
}

type taskEntry struct {
	taskHandler TaskHandler
	// eventFilters is a list of functions that are executed before a task is handled by the taskHandler. Only if all functions return 'true', the task will be handled
	eventFilters []func(keptnHandle IKeptn, event KeptnEvent) bool
}

func newTaskMap() *taskRegistry {
	return &taskRegistry{
		entries: make(map[string]taskEntry),
	}
}

func (t *taskRegistry) Contains(name string) (*taskEntry, bool) {
	t.RLock()
	defer t.RUnlock()
	if e, ok := t.entries[name]; ok {
		return &e, true
	} else if e, ok := t.entries["*"]; ok { // check if we have registered a wildcard handler
		return &e, true
	}
	return nil, false
}

func (t *taskRegistry) Add(name string, entry taskEntry) {
	t.Lock()
	defer t.Unlock()
	t.entries[name] = entry
}

func (t *taskRegistry) Get(name string) *taskEntry {
	t.RLock()
	defer t.RUnlock()
	entry := t.entries[name]
	return &entry
}
