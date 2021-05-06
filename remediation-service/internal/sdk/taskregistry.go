package sdk

import "sync"

type TaskRegistry struct {
	sync.RWMutex
	Entries map[string]TaskEntry
}

type TaskEntry struct {
	TaskHandler TaskHandler
	Context     Context
}

func NewTasksMap() TaskRegistry {
	return TaskRegistry{
		Entries: make(map[string]TaskEntry),
	}
}

func (t *TaskRegistry) Contains(name string) (*TaskEntry, bool) {
	t.Lock()
	defer t.Unlock()
	if e, ok := t.Entries[name]; ok {
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
