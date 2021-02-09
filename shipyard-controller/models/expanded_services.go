package models

// ExpandedServices expanded services
//
// swagger:model ExpandedProjects
type ExpandedServices struct {

	// Pointer to next page, base64 encoded
	NextPageKey string `json:"nextPageKey,omitempty"`

	// Size of returned page
	PageSize float64 `json:"pageSize,omitempty"`

	// projects
	Services []*ExpandedService `json:"services"`

	// Total number of projects
	TotalCount float64 `json:"totalCount,omitempty"`
}
