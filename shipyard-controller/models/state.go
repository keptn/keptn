package models

const (
	SequenceTriggeredState = "triggered"
	SequenceFinished       = "finished"
	TimedOut               = "timedOut"
)

type GetSequenceStateParams struct {
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

	/*Sequence name
	  In: query
	*/
	Name string `form:"name" json:"name"`

	/*Sequence status
	  In: query
	*/
	State string `form:"state" json:"state"`

	/*From time to fetch sequence states
	  In: query
	*/
	FromTime string `form:"fromTime" json:"fromTime"`

	/*Before time to fetch sequence states
	  In: query
	*/
	BeforeTime string `form:"beforeTime" json:"beforeTime"`

	/** Keptn context
	  In: query
	*/
	KeptnContext string `form:"keptnContext" json:"keptnContext"`
}

type StateFilter struct {
	GetSequenceStateParams
}

type SequenceStateEvaluation struct {
	Result string  `json:"result" bson:"result"`
	Score  float64 `json:"score" bson:"score"`
}

type SequenceStateEvent struct {
	Type string `json:"type" bson:"type"`
	ID   string `json:"id" bson:"id"`
	Time string `json:"time" bson:"time"`
}

type SequenceStateStage struct {
	Name              string                   `json:"name" bson:"name"`
	Image             string                   `json:"image,omitempty" bson:"image"`
	LatestEvaluation  *SequenceStateEvaluation `json:"latestEvaluation,omitempty" bson:"latestEvaluation"`
	LatestEvent       *SequenceStateEvent      `json:"latestEvent,omitempty" bson:"latestEvent"`
	LatestFailedEvent *SequenceStateEvent      `json:"latestFailedEvent,omitempty" bson:"latestFailedEvent"`
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
