package models

type PaginationParams struct {
	// NextPageKey indicates at which index the result should start
	NextPageKey int64 `json:"nextPageKey,omitempty"`

	// PageSize is the maximum size of returned page
	PageSize int64 `json:"pageSize,omitempty"`
}

type PaginationResult struct {
	// NextPageKey is the offset to the next page
	NextPageKey int64 `json:"nextPageKey,omitempty"`

	// PageSize is the actual size of returned page
	PageSize int64 `json:"pageSize,omitempty"`

	// Total number of matching entries
	TotalCount int64 `json:"totalCount,omitempty"`
}
