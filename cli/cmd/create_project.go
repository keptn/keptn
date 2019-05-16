package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/keptn/keptn/cli/utils/websockethelper"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type projectData struct {
	Project  string      `json:"project"`
	Registry interface{} `json:"registry"`
	Stages   interface{} `json:"stages"`
}

type myCloudEvent struct {
	contenttype string
	data        string
}

// crprojectCmd represents the project command
var crprojectCmd = &cobra.Command{
	Use:   "project project_name shipyard_file",
	Short: "Creates a new project.",
	Long: `Creates a new project with the provided name and shipyard file. 
The shipyard file describes the used stages.

Example:
	keptn create project sockshop shipyard.yml`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		_, _, err := credentialmanager.GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		if len(args) != 2 {
			cmd.SilenceUsage = false
			return errors.New("Requires project_name and shipyard_file")
		}
		if !utils.ValidateK8sName(args[0]) {
			errorMsg := "Project name includes invalid characters or is not well-formed.\n"
			errorMsg += "keptn relies on Helm charts and thus these conventions have to be followed: "
			errorMsg += "start with a lower case letter, then lower case letters, dash and numbers are allowed.\n"
			errorMsg += "You can find the guidelines here: https://github.com/helm/helm/blob/master/docs/chart_best_practices/conventions.md#chart-names\n"
			errorMsg += "Please update project name and try again."
			return errors.New(errorMsg)
		}
		if _, err := os.Stat(args[1]); os.IsNotExist(err) {
			return fmt.Errorf("Cannot find file %s", args[1])
		}
		_, err = utils.ReadFile(args[1])
		if err != nil {
			return err
		}
		testPrjData := projectData{}
		err = parseShipYard(&testPrjData, args[1])
		if err != nil {
			return fmt.Errorf("Invalid shipyard file")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}
		utils.PrintLog("Starting to create a project", utils.InfoLevel)

		prjData := projectData{}
		prjData.Project = args[0]

		parseShipYard(&prjData, args[1])

		source, _ := url.Parse("https://github.com/keptn/keptn/cli#createproject")

		contentType := "application/json"
		event := cloudevents.Event{
			Context: cloudevents.EventContextV02{
				ID:          uuid.New().String(),
				Type:        "create.project",
				Source:      types.URLRef{URL: *source},
				ContentType: &contentType,
			}.AsV02(),
			Data: prjData,
		}

		projectURL := endPoint
		projectURL.Path = "project"

		utils.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), utils.DebugLevel)

		if !mocking {
			responseCE, err := utils.Send(projectURL, event, apiToken)
			if err != nil {
				fmt.Println("Create project was unsuccessful")
				return err
			}

			// check for responseCE to include token
			if responseCE == nil {
				utils.PrintLog("Response CE is nil", utils.ErrorLevel)
				return nil
			}
			if responseCE.Data != nil {
				return websockethelper.PrintWSContent(responseCE)
			}
		} else {
			fmt.Println("Skipping create project due to mocking flag set to true")
		}
		return nil
	},
}

func parseShipYard(prjData *projectData, yamlFile string) error {
	data, err := utils.ReadFile(yamlFile)
	if err != nil {
		return err
	}

	var shipyardContent map[string]interface{}
	dec := yaml.NewDecoder(strings.NewReader(data))
	dec.Decode(&shipyardContent)

	if err != nil {
		return errors.New("Invalid shipyard file")
	}
	prjData.Registry = utils.Convert(shipyardContent["registry"])
	prjData.Stages = utils.Convert(shipyardContent["stages"])
	return nil
}

func init() {
	createCmd.AddCommand(crprojectCmd)
}
