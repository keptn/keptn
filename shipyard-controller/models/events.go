package models

// Events events
// swagger:model Events
type Events struct {

	// events
	Events []*Event `json:"events"`

	// Pointer to next page, base64 encoded
	NextPageKey string `json:"nextPageKey,omitempty"`

	// Size of returned page
	PageSize float64 `json:"pageSize,omitempty"`

	// Total number of events
	TotalCount float64 `json:"totalCount,omitempty"`
}
