package models

import (
	"encoding/json"
	"github.com/keptn/go-utils/pkg/api/models"
)

// Events events
// swagger:model Events
type Events struct {

	// events
	Events []*models.KeptnContextExtendedCE `json:"events"`

	// Pointer to next page, base64 encoded
	NextPageKey string `json:"nextPageKey,omitempty"`

	// Size of returned page
	PageSize float64 `json:"pageSize,omitempty"`

	// Total number of events
	TotalCount float64 `json:"totalCount,omitempty"`
}

type GetRootEventParams struct {
	Project     string `json:"project"`
	NextPageKey int64  `form:"nextPageKey" json:"nextPageKey"`
	PageSize    int64  `form:"pageSize" json:"pageSize"`
}

type GetEventsResult struct {
	// Pointer to next page
	NextPageKey int64 `json:"nextPageKey,omitempty"`

	// Size of returned page
	PageSize int64 `json:"pageSize,omitempty"`

	// Total number of logs
	TotalCount int64 `json:"totalCount,omitempty"`

	// Events
	Events []models.KeptnContextExtendedCE `json:"events"`
}

//ConvertToEvent returns an instance of models.Event, based on the provided input struct
func ConvertToEvent(in interface{}) (*models.KeptnContextExtendedCE, error) {
	bytes, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	result := &models.KeptnContextExtendedCE{}
	if err := json.Unmarshal(bytes, result); err != nil {
		return nil, err
	}
	return result, nil
}
