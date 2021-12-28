package models

import (
	"encoding/base64"
	"github.com/keptn/keptn/resource-service/common"
	"strings"
)

type ResourceContent string

func (rc ResourceContent) Validate() error {
	_, err := base64.StdEncoding.DecodeString(string(rc))
	if err != nil {
		return common.ErrResourceNotBase64Encoded
	}
	return nil
}

// Resource resource
//
// swagger:model Resource
type Resource struct {

	// Resource content - must be base64 encoded
	ResourceContent ResourceContent `json:"resourceContent,omitempty"`

	// Resource URI in URL-encoded format
	// Required: true
	ResourceURI string `json:"resourceURI"`
}

func (r Resource) Validate() error {
	if err := r.ResourceContent.Validate(); err != nil {
		return common.ErrResourceNotBase64Encoded
	}
	if err := validateResourceURI(r.ResourceURI); err != nil {
		return err
	}
	return nil
}

type GetResourcesQuery struct {
	GitCommitID string  `json:"gitCommitID,omitEmpty" form:"gitCommitID"`
	NextPageKey string  `json:"nextPageKey,omitempty" form:"nextPageKey"`
	PageSize    float64 `json:"pageSize,omitempty" form:"pageSize"`
}

type GetResourcesParams struct {
	Project
	Stage   *Stage
	Service *Service
	GetResourcesQuery
}

func (p GetResourcesParams) Validate() error {
	if err := p.Project.Validate(); err != nil {
		return err
	}
	if p.Stage != nil {
		if err := p.Stage.Validate(); err != nil {
			return err
		}
	}
	if p.Service != nil {
		if err := p.Service.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type GetResourceQuery struct {
	GitCommitID string `json:"gitCommitID,omitEmpty" form:"gitCommitID"`
}

type GetResourceParams struct {
	Project
	Stage       *Stage
	Service     *Service
	ResourceURI string
	GetResourceQuery
}

func (p GetResourceParams) Validate() error {
	if err := p.Project.Validate(); err != nil {
		return err
	}
	if p.Stage != nil {
		if err := p.Stage.Validate(); err != nil {
			return err
		}
	}
	if p.Service != nil {
		if err := p.Service.Validate(); err != nil {
			return err
		}
	}
	if err := validateResourceURI(p.ResourceURI); err != nil {
		return err
	}
	return nil
}

type DeleteResourceParams struct {
	Project
	Stage       *Stage
	Service     *Service
	ResourceURI string
}

func (p DeleteResourceParams) Validate() error {
	if err := p.Project.Validate(); err != nil {
		return err
	}
	if p.Stage != nil {
		if err := p.Stage.Validate(); err != nil {
			return err
		}
	}
	if p.Service != nil {
		if err := p.Service.Validate(); err != nil {
			return err
		}
	}
	if err := validateResourceURI(p.ResourceURI); err != nil {
		return err
	}
	return nil
}

type CreateResourceParams struct {
	Project
	Stage   *Stage
	Service *Service
	Resource
}

type UpdateResourcePayload struct {
	ResourceContent ResourceContent `json:"resourceContent"`
}

type UpdateResourceParams struct {
	Project
	Stage       *Stage
	Service     *Service
	ResourceURI string
	UpdateResourcePayload
}

func (p UpdateResourceParams) Validate() error {
	if err := p.Project.Validate(); err != nil {
		return err
	}
	if p.Stage != nil {
		if err := p.Stage.Validate(); err != nil {
			return err
		}
	}
	if p.Service != nil {
		if err := p.Service.Validate(); err != nil {
			return err
		}
	}
	if err := validateResourceURI(p.ResourceURI); err != nil {
		return err
	}
	if err := p.ResourceContent.Validate(); err != nil {
		return err
	}
	return nil
}

type CreateResourcesPayload struct {
	Resources []Resource `json:"resources"`
}

type CreateResourcesParams struct {
	Project
	Stage   *Stage
	Service *Service
	CreateResourcesPayload
}

func (p CreateResourcesParams) Validate() error {
	if err := p.Project.Validate(); err != nil {
		return err
	}
	if p.Stage != nil {
		if err := p.Stage.Validate(); err != nil {
			return err
		}
	}
	if p.Service != nil {
		if err := p.Service.Validate(); err != nil {
			return err
		}
	}
	for _, res := range p.Resources {
		if err := res.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type UpdateResourcesPayload struct {
	Resources []Resource `json:"resources"`
}

type UpdateResourcesParams struct {
	Project
	Stage   *Stage
	Service *Service
	UpdateResourcesPayload
}

func (p UpdateResourcesParams) Validate() error {
	if err := p.Project.Validate(); err != nil {
		return err
	}
	if p.Stage != nil {
		if err := p.Stage.Validate(); err != nil {
			return err
		}
	}
	if p.Service != nil {
		if err := p.Service.Validate(); err != nil {
			return err
		}
	}
	for _, res := range p.Resources {
		if err := res.Validate(); err != nil {
			return err
		}
	}
	return nil
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
	Resources []GetResourceResponse `json:"resources"`

	// Total number of resources
	TotalCount float64 `json:"totalCount,omitempty"`
}

// GetResourceResponse resources
//
// swagger:model GetResourceResponse
type GetResourceResponse struct {
	Resource
	Metadata Version
}

func validateResourceURI(uri string) error {
	if strings.Contains(uri, "~") || strings.Contains(uri, "..") {
		return common.ErrResourceInvalidResourceURI
	}
	return nil
}
