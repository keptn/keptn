package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/knative/pkg/cloudevents"
	"github.com/spf13/cobra"
)

type projectData struct {
	Project string        `json:"project"`
	Stages  []utils.Stage `json:"stages"`
}

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Creates a new repository in the GitHub organization and initializes the repository with helm charts.",
	Long: `Creates a new repository in the GitHub organization and initializes the repository with helm charts. 
	Therfore, the provided name and the description for the stages (specified in the provided yml file)
	is used. Usage of \"create project\":

keptn create project sockshop shipyard.yml`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("Requires exact two args specifying the project name and the description for the stages")
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
			return fmt.Errorf("No stages defined")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
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
		endPoint, apiToken, err := credentialmanager.GetCreds()
		if err != nil || endPoint == "" {
			utils.Info.Printf("create project called without beeing authenticated.")
			return errors.New("This command requires to be authenticated. See \"keptn auth\" for details")
		}
		req, err := builder.Build(endPoint+"project", prjData)
		if err != nil {
			return err
		}

		err = utils.Send(req, apiToken)
		if err != nil {
			utils.Error.Printf("create project command was unsuccessful. Details: %v", err)
			return err
		}
		fmt.Printf("Successfully created project %v on Github\n", prjData.Project)
		return nil
	},
}

func init() {
	createCmd.AddCommand(projectCmd)
}
