package operations

import (
	"errors"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
)

type CreateServiceParams struct {
	// name
	ServiceName *string `json:"serviceName"`
}

type CreateServiceResponse struct {
}

type DeleteServiceResponse struct {
	Message string `json:"message"`
}

func (params *CreateServiceParams) Validate() error {
	if params.ServiceName == nil || *params.ServiceName == "" {
		return errors.New("Must provide a service name")
	}
	if !keptncommon.ValidateUnixDirectoryName(*params.ServiceName) {
		return errors.New("Service name contains special character(s). " +
			"The service name has to be a valid Unix directory name. For details see " +
			"https://www.cyberciti.biz/faq/linuxunix-rules-for-naming-file-and-directory-names/")
	}
	return nil
}

type GetServiceParams struct {

	//Pointer to the next set of items
	NextPageKey *string `form:"nextPageKey"`

	//The number of items to return
	PageSize *int64 `form:"pageSize"`
}
