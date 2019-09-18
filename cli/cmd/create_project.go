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
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/keptn/keptn/cli/utils/websockethelper"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type createProjectCmdParams struct {
	GitUser   *string
	GitToken  *string
	RemoteURL *string
}

var createProjectParams *createProjectCmdParams

// crprojectCmd represents the project command
var crprojectCmd = &cobra.Command{
	Use:   "project <project_name> <shipyard_file> --git-user=<git_user> --git-token=<git_token> --git-remote-url=<git_remote_url>",
	Short: "Creates a new project.",
	Long: `Creates a new project with the provided name and shipyard file. 
The shipyard file describes the used stages. Furthermore, for these stages the shipyard file 
describes the used deployment and test strategies.

Example:
	keptn create project project_name shipyard_file.yml`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		_, _, err := credentialmanager.GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		if len(args) != 2 {
			cmd.SilenceUsage = false
			return errors.New("Requires project_name and shipyard_file.yml")
		}
		if !utils.ValidateK8sName(args[0]) {
			errorMsg := "Project name includes invalid characters or is not well-formed.\n"
			errorMsg += "keptn relies on Helm charts and thus these conventions have to be followed: "
			errorMsg += "start with a lower case letter, then lower case letters, dash and numbers are allowed.\n"
			errorMsg += "You can find the guidelines here: https://github.com/helm/helm/blob/master/docs/chart_best_practices/conventions.md#chart-names\n"
			errorMsg += "Please update project name and try again."
			return errors.New(errorMsg)
		}
		if _, err := os.Stat(keptnutils.ExpandTilde(args[1])); os.IsNotExist(err) {
			return fmt.Errorf("Cannot find file %s", keptnutils.ExpandTilde(args[1]))
		}

		content, err := utils.ReadFile(keptnutils.ExpandTilde(args[1]))
		if err != nil {
			return err
		}

		if err := testParseShipyard(content); err != nil {
			return fmt.Errorf("Invalid shipyard file because parsing failed: %s", err.Error())
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		gitUser := true
		gitToken := true
		remoteURL := true

		if *createProjectParams.GitUser == "" {
			gitUser = false
		}
		if *createProjectParams.GitToken == "" {
			gitToken = false
		}
		if *createProjectParams.RemoteURL == "" {
			remoteURL = false
		}

		if gitUser == false && gitToken == false && remoteURL == false {
			return nil
		}

		if gitUser != true || gitToken != true || remoteURL != true {
			return errors.New("For configuring a Git upstream repository please specify Git user, user token, and remote URL of repository/project")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}
		logging.PrintLog("Starting to create a project", logging.InfoLevel)

		content, _ := utils.ReadFile(keptnutils.ExpandTilde(args[1]))
		prjData := keptnevents.ProjectCreateEventData{Project: args[0], Shipyard: base64.StdEncoding.EncodeToString([]byte(content))}

		if *createProjectParams.GitUser != "" && *createProjectParams.GitToken != "" && *createProjectParams.RemoteURL != "" {
			prjData.GitUser = *createProjectParams.GitUser
			prjData.GitToken = *createProjectParams.GitToken
			prjData.GitRemoteURL = *createProjectParams.RemoteURL
		}

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

		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		if !mocking {
			responseCE, err := utils.Send(projectURL, event, apiToken)
			if err != nil {
				fmt.Println("Create project was unsuccessful")
				return err
			}

			// check for responseCE to include token
			if responseCE == nil {
				logging.PrintLog("Response CE is nil", logging.QuietLevel)
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

func testParseShipyard(shipyardContent string) error {
	shipyard := models.Shipyard{}
	return yaml.Unmarshal([]byte(shipyardContent), &shipyard)
}

func init() {
	createCmd.AddCommand(crprojectCmd)
	createProjectParams = &createProjectCmdParams{}
	createProjectParams.GitUser = crprojectCmd.Flags().StringP("git-user", "u", "", "The git user of the upstream target")
	createProjectParams.GitToken = crprojectCmd.Flags().StringP("git-token", "t", "", "The git token of the git user")
	createProjectParams.RemoteURL = crprojectCmd.Flags().StringP("git-remote-url", "r", "", "The remote url of the upstream target")
}
