package sdk

type TaskRegistry struct {
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

func (tm *TaskRegistry) Contains(name string) (*TaskEntry, bool) {
	if e, ok := tm.Entries[name]; ok {
		return &e, true
	}
	return nil, false
}

func (tm *TaskRegistry) Add(name string, entry TaskEntry) {
	tm.Entries[name] = entry
}

func (tm *TaskRegistry) Get(name string) *TaskEntry {
	entry := tm.Entries[name]
	return &entry
}
