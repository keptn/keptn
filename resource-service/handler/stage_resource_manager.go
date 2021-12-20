package handler

type IStageResourceManager interface {
}

type StageResourceManager struct {
}

func NewStageResourceManager() *StageResourceManager {
	stageResourceManager := &StageResourceManager{}
	return stageResourceManager
}
