package handler

type IStageManager interface {
}

type StageManager struct {
}

func NewStageManager() *StageManager {
	stageManager := &StageManager{}
	return stageManager
}
