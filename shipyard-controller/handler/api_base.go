package handler

import (
	"fmt"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/keptn/shipyard-controller/common"
)

type apiBase struct {
	projectAPI  *keptnapi.ProjectHandler
	stagesAPI   *keptnapi.StageHandler
	servicesAPI *keptnapi.ServiceHandler
	resourceAPI *keptnapi.ResourceHandler
	secretStore common.SecretStore
	logger      keptncommon.LoggerInterface
}

func newAPIBase() (*apiBase, error) {
	csEndpoint, err := keptncommon.GetServiceEndpoint("CONFIGURATION_SERVICE")
	if err != nil {
		return nil, fmt.Errorf("could not get configuration-service URL: %s", err.Error())
	}
	secretStore, err := common.NewK8sSecretStore()
	if err != nil {
		return nil, fmt.Errorf("could not initilize secret store: " + err.Error())
	}

	return &apiBase{
		projectAPI:  keptnapi.NewProjectHandler(csEndpoint.String()),
		stagesAPI:   keptnapi.NewStageHandler(csEndpoint.String()),
		servicesAPI: keptnapi.NewServiceHandler(csEndpoint.String()),
		resourceAPI: keptnapi.NewResourceHandler(csEndpoint.String()),
		secretStore: secretStore,
		logger:      keptncommon.NewLogger("", "", "shipyard-controller"),
	}, nil
}
