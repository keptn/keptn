package models

// Resource resource
//
// swagger:model Resource
type Resource struct {

	// metadata
	Metadata *Version `json:"metadata,omitempty"`

	// Resource content
	ResourceContent string `json:"resourceContent,omitempty"`

	// Resource URI in URL-encoded format
	// Required: true
	ResourceURI string `json:"resourceURI"`
}

type CreateResourceParams struct {
	// Resource content
	ResourceContent string `json:"resourceContent,omitempty"`

	// Resource URI in URL-encoded format
	// Required: true
	ResourceURI string `json:"resourceURI"`
}

type UpdateResourceParams CreateResourceParams

type CreateResourcesParams struct {
	Resources []CreateResourceParams `json:"resources"`
}

type UpdateResourcesParams struct {
	Resources []UpdateResourceParams `json:"resources"`
}

// GetResourcesResponse resources
//
// swagger:model GetResourcesResponse
type GetResourcesResponse struct {

	// Pointer to next page, base64 encoded
	NextPageKey string `json:"nextPageKey,omitempty"`

	// Size of returned page
	PageSize float64 `json:"pageSize,omitempty"`

	// resources
	Resources []*Resource `json:"resources"`

	// Total number of resources
	TotalCount float64 `json:"totalCount,omitempty"`
}

// GetResourceResponse resources
//
// swagger:model GetResourceResponse
type GetResourceResponse Resource
