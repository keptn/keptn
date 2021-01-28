package operations

// Event godoc
type Event struct {
	Contenttype    string      `json:"contenttype,omitempty"`
	Data           KeptnBase   `json:"data"`
	Extensions     interface{} `json:"extensions,omitempty"`
	ID             string      `json:"id,omitempty"`
	Shkeptncontext string      `json:"shkeptncontext,omitempty"`
	Source         string      `json:"source"`
	Specversion    string      `json:"specversion,omitempty"`
	Time           string      `json:"time,omitempty"`
	Triggeredid    string      `json:"triggeredid,omitempty"`
	Type           string      `json:"type"`
}

// KeptnBase godoc
type KeptnBase struct {
	Project string `json:"project"`
	Service string `json:"service"`
}
