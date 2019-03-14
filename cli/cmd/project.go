package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/keptn/keptn/cli/utils/websockethelper"
	"github.com/knative/pkg/cloudevents"
	"github.com/spf13/cobra"
)

type projectData struct {
	Project string        `json:"project"`
	Stages  []utils.Stage `json:"stages"`
}

// projectCmd represents the project command
var projectCmd = &cobra.Command{
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
			return errors.New("Requires project_name and shipyard_file\n")
		}
		if _, err := os.Stat(args[1]); os.IsNotExist(err) {
			return fmt.Errorf("Cannot find file %s", args[1])
		}
		data, err := ioutil.ReadFile(args[1])
		if err != nil {
			return err
		}
		stages := utils.UnmarshalStages(data)
		if len(stages) == 0 {
			return fmt.Errorf("No stages defined in provided shipyard file")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		fmt.Println("Starting to create a project")

		prjData := projectData{}
		prjData.Project = args[0]
		data, err := ioutil.ReadFile(args[1])
		if err != nil {
			return err
		}
		prjData.Stages = utils.UnmarshalStages(data)

		builder := cloudevents.Builder{
			Source:    "https://github.com/keptn/keptn/cli#createproject",
			EventType: "create.project",
			Encoding:  cloudevents.StructuredV01,
		}

		projectURL := endPoint
		projectURL.Path = "project"

		req, err := builder.Build(projectURL.String(), prjData)
		if err != nil {
			return err
		}

		var desc = new(utils.WebsocketDescription)
		resp, err := utils.Send(req, apiToken, desc)

		if err != nil {
			fmt.Println("Create project was unsuccessful")
			return err
		}
		if resp.StatusCode != 200 {
			fmt.Println("Create project was unsuccessful")
			return errors.New(resp.Status)
		}

		if desc.Token != "" {
			ws, err := websockethelper.OpenWS(desc.Token)
			if err != nil {
				return err
			}
			return websockethelper.PrintWSContent(ws)
		}

		fmt.Printf("Successfully created project %v on Github\n", prjData.Project)
		return nil
	},
}

func init() {
	createCmd.AddCommand(projectCmd)
}
