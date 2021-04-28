package models

type GetStateParams struct {
	/*Pointer to the next set of items
	  In: query
	*/
	NextPageKey int64 `form:"nextPageKey" json:"nextPageKey"`
	/*The number of items to return
	  Maximum: 50
	  Minimum: 1
	  In: query
	  Default: 20
	*/
	PageSize int64 `form:"pageSize" json:"pageSize"`
	/*Project name
	  In: query
	*/
	Project string `form:"project" json:"project"`
}

type StateFilter struct {
	GetStateParams
	Shkeptncontext string
	Name           string
}

type SequenceStateEvaluation struct {
	Result string  `json:"result" bson:"result"`
	Score  float64 `json:"score" bson:"score"`
}

type SequenceStateEvent struct {
	Type   string `json:"type" bson:"type"`
	ID     string `json:"id" bson:"id"`
	Time   string `json:"time" bson:"time"`
	Result string `json:"result" bson:"result"`
}

type SequenceStateStage struct {
	Name             string                  `json:"name" bson:"name"`
	Image            string                  `json:"image" bson:"image"`
	LatestEvaluation SequenceStateEvaluation `json:"latestEvaluation" bson:"latestEvaluation"`
	LatestEvent      SequenceStateEvent      `json:"latestEvent" bson:"latestEvent"`
}

type SequenceState struct {
	Name           string               `json:"name" bson:"name"`
	Service        string               `json:"service" bson:"service"`
	Project        string               `json:"project" bson:"project"`
	Time           string               `json:"time" bson:"time"`
	Shkeptncontext string               `json:"shkeptncontext" bson:"shkeptncontext"`
	State          string               `json:"state" bson:"state"`
	Stages         []SequenceStateStage `json:"stages" bson:"stages"`
}

type SequenceStates struct {
	States []SequenceState `json:"states"`
	// Pointer to next page
	NextPageKey int64 `json:"nextPageKey,omitempty"`

	// Size of returned page
	PageSize int64 `json:"pageSize,omitempty"`

	// Total number of events
	TotalCount int64 `json:"totalCount,omitempty"`
}
