package sdk

import (
	"sync"
)

type TaskRegistry struct {
	sync.RWMutex
	Entries map[string]TaskEntry
}

type TaskEntry struct {
	TaskHandler TaskHandler
	// EventFilters is a list of functions that are executed before a task is handled by the TaskHandler. Only if all functions return 'true', the task will be handled
	EventFilters []func(keptnHandle IKeptn, event KeptnEvent) bool
}

func NewTasksMap() *TaskRegistry {
	return &TaskRegistry{
		Entries: make(map[string]TaskEntry),
	}
}

func (t *TaskRegistry) Contains(name string) (*TaskEntry, bool) {
	t.Lock()
	defer t.Unlock()
	if e, ok := t.Entries[name]; ok {
		return &e, true
	} else if e, ok := t.Entries["*"]; ok { // check if we have registered a wildcard handler
		return &e, true
	}
	return nil, false
}

func (t *TaskRegistry) Add(name string, entry TaskEntry) {
	t.Lock()
	defer t.Unlock()
	t.Entries[name] = entry
}

func (t *TaskRegistry) Get(name string) *TaskEntry {
	t.Lock()
	defer t.Unlock()
	entry := t.Entries[name]
	return &entry
}
