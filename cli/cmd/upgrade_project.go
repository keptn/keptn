package cmd

import (
	"bufio"
	"errors"
	"fmt"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	keptn "github.com/keptn/go-utils/pkg/lib"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"log"
	"net/url"
	"os"
	"strings"
)

type upgradeProjectCmdParams struct {
	Shipyard    bool
	DryRun      bool
	FromVersion *string
	ToVersion   *string
	AutoConfirm bool
}

var upgradeProjectParams *upgradeProjectCmdParams

var supportedFromVersions = []string{"0.1", "0.1.0"}
var supportedToVersions = []string{"0.2", "0.2.0"}

const defaultFromVersion = "0.1"
const defaultToVersion = "0.2"

// upgradeProjectCmd represents the project command
var upgradeProjectCmd = &cobra.Command{
	Use:   "project PROJECTNAME --shipyard --fromVersion=CURRENT_SHIPYARD_VERSION --toVersion=TARGET_SHIPYARD_VERSION",
	Short: "Upgrades an existing Keptn project",
	Long: `Upgrades an existing Keptn project with the provided name. 

This command will upgrade the shipyard of the project to the specified version

By executing the update project command, Keptn will fetch the current shipyard.yaml file of the project and convert it to the version specified in the 'toVersion'' flag.

For more information about upgrading projects, go to [Manage Keptn](https://keptn.sh/docs/` + keptnReleaseDocsURL + `/operate/upgrade)
`,
	Example:      `keptn upgrade project PROJECTNAME --shipyard --fromVersion=0.1.0 --toVersion=0.2.0`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		_, _, err := credentialmanager.NewCredentialManager(false).GetCreds(namespace)
		if err != nil {
			return errors.New(authErrorMsg)
		}

		if len(args) != 1 {
			cmd.SilenceUsage = false
			return errors.New("required argument PROJECTNAME not set")
		}

		if upgradeProjectParams.FromVersion == nil || *upgradeProjectParams.FromVersion == "" {
			*upgradeProjectParams.FromVersion = defaultFromVersion
		} else if err := checkFromVersion(upgradeProjectParams.FromVersion); err != nil {
			return err
		}

		if upgradeProjectParams.ToVersion == nil || *upgradeProjectParams.ToVersion == "" {
			*upgradeProjectParams.ToVersion = defaultToVersion
		} else if err := checkToVersion(upgradeProjectParams.ToVersion); err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		var endPoint url.URL
		var apiToken string
		var err error
		if !mocking {
			endPoint, apiToken, err = credentialmanager.NewCredentialManager(false).GetCreds(namespace)
		} else {
			endPointPtr, _ := url.Parse(os.Getenv("MOCK_SERVER"))
			endPoint = *endPointPtr
			apiToken = ""
		}
		if endPointErr := checkEndPointStatus(endPoint.String()); endPointErr != nil {
			return fmt.Errorf("Error connecting to server: %s"+endPointErrorReasons,
				endPointErr)
		}

		resourceHandler := apiutils.NewAuthenticatedResourceHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
		shipyardResource, err := resourceHandler.GetProjectResource(projectName, "shipyard.yaml")
		if err != nil {
			return fmt.Errorf("Error while retrieving shipyard.yaml for project %s: %s:", *newArtifact.Project, err.Error())
		}

		// first, check if the shipyard already has been upgraded
		alreadyUpgraded, err := isShipyardUpgraded(shipyardResource)
		if err != nil {
			return fmt.Errorf("could not check if shipyard of project %s is already up to date: %s", projectName, err.Error())
		}
		if alreadyUpgraded {
			logging.PrintLog("Shipyard of project "+projectName+" has already been upgraded to version 0.2", logging.InfoLevel)
			return nil
		}

		shipyard := &keptn.Shipyard{}
		if err := yaml.Unmarshal([]byte(shipyardResource.ResourceContent), shipyard); err != nil {
			return fmt.Errorf("error while decoding shipyard.yaml for project %s: %s", *newArtifact.Project, err.Error())
		}

		// check if there are any stages in the old shipyard.
		// Having a shipyard with no stage should not happen, so this would mean that something has gone wrong when unmarshalling into the struct.
		// in this case, the upgrade is cancelled to avoid deleting data
		if len(shipyard.Stages) == 0 {
			logging.PrintLog("Current shipyard.yaml of project "+projectName+" does not contain any stages. Will not proceed with upgrade", logging.InfoLevel)
			return nil
		}

		upgradedShipyard := transformShipyard(shipyard)
		marshalledUpgradedShipyard, err := yaml.Marshal(upgradedShipyard)
		if err != nil {
			return fmt.Errorf("could not marshal upgraded shipyard into string: %s", err.Error())
		}

		logging.PrintLog("Shipyard of project "+projectName+":", logging.InfoLevel)
		logging.PrintLog("-----------------------", logging.InfoLevel)
		logging.PrintLog(string(shipyardResource.ResourceContent), logging.InfoLevel)

		logging.PrintLog("Shipyard converted into version 0.2:", logging.InfoLevel)
		logging.PrintLog("-----------------------", logging.InfoLevel)
		logging.PrintLog(string(marshalledUpgradedShipyard), logging.InfoLevel)

		if upgradeProjectParams.DryRun {
			return nil
		}

		if err := confirmShipyardUpgrade(); err != nil {
			return err
		}

		shipyardName := "shipyard.yaml"
		upgradedShipyardResource := &apimodels.Resource{
			ResourceContent: string(marshalledUpgradedShipyard),
			ResourceURI:     &shipyardName,
		}
		if _, err := resourceHandler.UpdateProjectResource(projectName, upgradedShipyardResource); err != nil {
			return fmt.Errorf("could not update shipyard resource: %s", err.Error())
		}
		logging.PrintLog("Shipyard of project "+projectName+" has been upgraded successfully!", logging.InfoLevel)

		return nil
	},
}

func isShipyardUpgraded(resource *apimodels.Resource) (bool, error) {
	v2Shipyard := &keptnv2.Shipyard{}

	if err := yaml.Unmarshal([]byte(resource.ResourceContent), v2Shipyard); err != nil {
		return false, err
	}

	if strings.Contains(v2Shipyard.ApiVersion, *upgradeProjectParams.ToVersion) {
		return true, nil
	}
	return false, nil
}

func confirmShipyardUpgrade() error {
	if upgradeProjectParams.AutoConfirm {
		return nil
	}
	logging.PrintLog("Do you want to continue with this? (y/n)", logging.InfoLevel)
	reader := bufio.NewReader(os.Stdin)
	in, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	in = strings.ToLower(strings.TrimSpace(in))
	if !(in == "y" || in == "yes") {
		err := errors.New("stopping installation")
		log.Fatal(err)
	}
	return nil
}

func transformShipyard(shipyard *keptn.Shipyard) *keptnv2.Shipyard {
	upgradedShipyard := &keptnv2.Shipyard{
		ApiVersion: "spec.keptn.sh/0.2.0",
		Kind:       "Shipyard",
		Spec: keptnv2.ShipyardSpec{
			Stages: []keptnv2.Stage{},
		},
	}

	bytes, _ := yaml.Marshal(shipyard)

	fmt.Println(string(bytes))

	for index, stage := range shipyard.Stages {

		passStrategy, warningStrategy := getApprovalStrategyForStage(index, shipyard)
		newStage := keptnv2.Stage{
			Name: stage.Name,
			Sequences: []keptnv2.Sequence{
				{
					Name:     "artifact-delivery",
					Triggers: getSequenceTriggerForStage(index, shipyard, "artifact-delivery"),
					Tasks: []keptnv2.Task{
						{
							Name: "deployment",
							Properties: map[string]string{
								"deploymentstrategy": stage.DeploymentStrategy,
							},
						},
						{
							Name: "test",
							Properties: map[string]string{
								"teststrategy": stage.TestStrategy,
							},
						},
						{
							Name: "evaluation",
						},
						{
							Name: "approval",
							Properties: map[string]string{
								"pass":    passStrategy,
								"warning": warningStrategy,
							},
						},
						{
							Name: "release",
						},
					},
				},
				// add a second artifact-delivery with "direct" deployment strategy
				{
					Name:     "artifact-delivery-direct",
					Triggers: getSequenceTriggerForStage(index, shipyard, "artifact-delivery-direct"),
					Tasks: []keptnv2.Task{
						{
							Name: "deployment",
							Properties: map[string]string{
								"deploymentstrategy": "direct",
							},
						},
						{
							Name: "test",
							Properties: map[string]string{
								"teststrategy": stage.TestStrategy,
							},
						},
						{
							Name: "evaluation",
						},
						{
							Name: "approval",
							Properties: map[string]string{
								"pass":    passStrategy,
								"warning": warningStrategy,
							},
						},
						{
							Name: "release",
						},
					},
				},
			},
		}
		upgradedShipyard.Spec.Stages = append(upgradedShipyard.Spec.Stages, newStage)
	}

	return upgradedShipyard
}

func getApprovalStrategyForStage(index int, shipyard *keptn.Shipyard) (string, string) {
	if shipyard.Stages[index].ApprovalStrategy == nil {
		return keptn.Automatic.String(), keptn.Automatic.String()
	}

	return shipyard.Stages[index].ApprovalStrategy.Pass.String(), shipyard.Stages[index].ApprovalStrategy.Warning.String()
}

func getSequenceTriggerForStage(index int, shipyard *keptn.Shipyard, sequenceName string) []string {
	if index == 0 {
		return []string{}
	}
	return []string{shipyard.Stages[index-1].Name + "." + sequenceName + ".finished"}
}

func checkFromVersion(fromVersion *string) error {
	for _, value := range supportedFromVersions {
		if value == *fromVersion {
			return nil
		}
	}
	return fmt.Errorf("invalid fromVersion %s. Please enter one of the following: %v", *fromVersion, supportedFromVersions)
}

func checkToVersion(toVersion *string) error {
	for _, value := range supportedToVersions {
		if value == *toVersion {
			return nil
		}
	}
	return fmt.Errorf("invalid toVersion %s. Please enter one of the following: %v", *toVersion, supportedToVersions)
}

func init() {
	upgraderCmd.AddCommand(upgradeProjectCmd)

	upgradeProjectParams = &upgradeProjectCmdParams{}

	upgradeProjectCmd.Flags().BoolVarP(&upgradeProjectParams.Shipyard, "shipyard", "", false, "Upgrade the shipyard file of the project")
	upgradeProjectCmd.Flags().BoolVarP(&upgradeProjectParams.DryRun, "dry-run", "", false, "Output the upgraded shipyard but don't upload it to the project")
	upgradeProjectParams.FromVersion = upgradeProjectCmd.Flags().StringP("fromVersion", "", "", "The current version of the shipyard")
	upgradeProjectParams.ToVersion = upgradeProjectCmd.Flags().StringP("toVersion", "", "", "The new target version of the shipyard")
	upgradeProjectCmd.Flags().BoolVarP(&upgradeProjectParams.AutoConfirm, "yes", "y", false, "Automatically confirm the upgrade of the shipyard")
	upgradeProjectCmd.MarkFlagRequired("shipyard")
}
