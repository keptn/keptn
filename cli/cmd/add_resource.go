package cmd

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/keptn/go-utils/pkg/models"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/spf13/cobra"
)

type addResourceCommandParameters struct {
	Project     *string
	Stage       *string
	Service     *string
	Resource    *string
	ResourceURI *string
}

type addResourceError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var addResourceCmdParams *addResourceCommandParameters

var addResourceCmd = &cobra.Command{
	Use:   "add-resource --project=PROJECT --stage=STAGE --service=SERVICE --resource=FILEPATH --resourceUri=FILEPATH",
	Short: "Adds a resource to a service within your project in the specified stage",
	Long: `Adds a resource to a service within your project in the specified stage
	
Example: 
	keptn add-resource --project=sockshop --stage=dev --service=carts --resource=./jmeter.jmx --resourceUri=jmeter/functional.jmx`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}
		*addResourceCmdParams.Resource = keptnutils.ExpandTilde(*addResourceCmdParams.Resource)

		if !fileExists(*addResourceCmdParams.Resource) {
			return errors.New("File " + *addResourceCmdParams.Resource + " not found in local file system")
		}
		resourceContent, err := ioutil.ReadFile(*addResourceCmdParams.Resource)
		if err != nil {
			return errors.New("File " + *addResourceCmdParams.Resource + " could not be read")
		}

		logging.PrintLog("Adding resource "+*addResourceCmdParams.Resource+" to service "+*addResourceCmdParams.Service+" in stage "+*addResourceCmdParams.Stage+" in project "+*addResourceCmdParams.Project, logging.InfoLevel)

		if *addResourceCmdParams.ResourceURI == "" {
			addResourceCmdParams.ResourceURI = addResourceCmdParams.Resource
		}
		resources := []*models.Resource{
			&models.Resource{
				ResourceContent: string(resourceContent),
				ResourceURI:     addResourceCmdParams.ResourceURI,
			},
		}

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			DialContext:     utils.ResolveXipIoWithContext,
		}
		client := &http.Client{Transport: tr}
		resourceHandler := keptnutils.NewAuthenticatedResourceHandler(endPoint.Host, apiToken, "x-token", client, "https")
		_, err = resourceHandler.CreateServiceResources(*addResourceCmdParams.Project, *addResourceCmdParams.Stage, *addResourceCmdParams.Service, resources)
		if err != nil {
			errorObj := &addResourceError{}
			err2 := json.Unmarshal([]byte(err.Error()), errorObj)
			if err2 != nil {
				return errors.New("Resource " + *addResourceCmdParams.Resource + " could not be uploaded: " + err.Error())
			}
			return errors.New("Resource " + *addResourceCmdParams.Resource + " could not be uploaded: " + errorObj.Message)
		}
		logging.PrintLog("Resource has been uploaded.", logging.InfoLevel)
		return nil
	},
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func init() {
	rootCmd.AddCommand(addResourceCmd)
	addResourceCmdParams = &addResourceCommandParameters{}
	addResourceCmdParams.Project = addResourceCmd.Flags().StringP("project", "p", "", "The name of the project")
	addResourceCmd.MarkFlagRequired("project")
	addResourceCmdParams.Stage = addResourceCmd.Flags().StringP("stage", "s", "", "The name of the stage")
	addResourceCmd.MarkFlagRequired("stage")
	addResourceCmdParams.Service = addResourceCmd.Flags().StringP("service", "", "", "The name of the service within the project")
	addResourceCmd.MarkFlagRequired("service")
	addResourceCmdParams.Resource = addResourceCmd.Flags().StringP("resource", "r", "", "Path pointing to the resource on your local file system")
	addResourceCmd.MarkFlagRequired("resource")
	addResourceCmdParams.ResourceURI = addResourceCmd.Flags().StringP("resourceUri", "", "", "Optional: Location where the resource should be stored within the config repo. If empty, The name of the resource will be the same as on your local file system")

}
