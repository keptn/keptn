package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	"github.com/keptn/go-utils/pkg/models"
	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/keptn/keptn/cli/utils/websockethelper"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// crprojectCmd represents the project command
var crprojectCmd = &cobra.Command{
	Use:   "project project_name shipyard_file",
	Short: "Creates a new project.",
	Long: `Creates a new project with the provided name and shipyard file. 
The shipyard file describes the used stages. Furthermore, for these stages the shipyard file 
describes the used deployment and test strategies.

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
		content, err := utils.ReadFile(args[1])
		if err != nil {
			return err
		}

		if err := testParseShipYard(content); err != nil {
			return fmt.Errorf("Invalid shipyard file because parsing failed: %s", err.Error())
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}
		utils.PrintLog("Starting to create a project", utils.InfoLevel)

		content, _ := utils.ReadFile(args[1])
		prjData := keptnevents.ProjectCreateEventData{Project: args[0], Shipyard: base64.StdEncoding.EncodeToString([]byte(content))}

		source, _ := url.Parse("https://github.com/keptn/keptn/cli#createproject")

		contentType := "application/json"
		event := cloudevents.Event{
			Context: cloudevents.EventContextV02{
				ID:          uuid.New().String(),
				Type:        keptnevents.InternalProjectCreateEventType,
				Source:      types.URLRef{URL: *source},
				ContentType: &contentType,
			}.AsV02(),
			Data: prjData,
		}

		projectURL := endPoint
		projectURL.Path = "v1/project"

		utils.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), utils.VerboseLevel)

		if !mocking {
			responseCE, err := utils.Send(projectURL, event, apiToken)
			if err != nil {
				fmt.Println("Create project was unsuccessful")
				return err
			}

			// check for responseCE to include token
			if responseCE == nil {
				utils.PrintLog("Response CE is nil", utils.QuietLevel)
				return nil
			}
			if responseCE.Data != nil {
				return websockethelper.PrintWSContentCEResponse(responseCE, endPoint)
			}
		} else {
			fmt.Println("Skipping create project due to mocking flag set to true")
		}
		return nil
	},
}

func testParseShipYard(shipyardContent string) error {
	shipyard := models.Shipyard{}
	return yaml.Unmarshal([]byte(shipyardContent), &shipyard)
}

func init() {
	createCmd.AddCommand(crprojectCmd)
}
