package db

type TaskSequenceRepo interface {
	GetTaskSequence(project, triggeredID string) (string, error)
}
