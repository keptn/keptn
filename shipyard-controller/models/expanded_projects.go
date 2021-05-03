package models

// ExpandedProjects expanded projects
//
// swagger:model ExpandedProjects
type ExpandedProjects struct {

	// Pointer to next page, base64 encoded
	NextPageKey string `json:"nextPageKey,omitempty"`

	// Size of returned page
	PageSize float64 `json:"pageSize,omitempty"`

	// projects
	Projects []*ExpandedProject `json:"projects"`

	// Total number of projects
	TotalCount float64 `json:"totalCount,omitempty"`
}
