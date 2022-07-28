package sequence_execution

import "github.com/keptn/keptn/shipyard-controller/models"

// ModelTransformer is an interface that defines functions for transforming between the internal representation
// of a sequence execution and the model structure outside the db package
type ModelTransformer interface {
	TransformToDBModel(execution models.SequenceExecution) interface{}
	TransformEventToDBModel(event models.TaskEvent) interface{}
	TransformToSequenceExecution(dbItem interface{}) (*models.SequenceExecution, error)
}
