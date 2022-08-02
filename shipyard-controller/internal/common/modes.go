package common

type SDMode int

const (
	// SDModeRW specifies that Sequence Dispatcher has both write and read roles
	SDModeRW SDMode = iota
	// SDModeW set Sequence Dispatcher to act in write only mode, this is needed with multiple replicas
	SDModeW
)
