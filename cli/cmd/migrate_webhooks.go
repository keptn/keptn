package cmd

import (
	"fmt"

	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/internal"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/webhook-service/lib"
	"github.com/spf13/cobra"
)

type migrateWebhooksCmdParams struct {
	ProjectName *string
	DryRun      *bool
	AcceptAll   bool
}

var migrateWebhooksParams *migrateWebhooksCmdParams
var ErrResourceNotFound = fmt.Errorf("Resource not found")

var migrateWebhooksCmd = &cobra.Command{
	Use:     "migrate-webhooks",
	Short:   `Migrates the webhook configurations from version v1alpha1 to v1beta1`,
	Example: `keptn migrate-webhooks [--project=PROJECTMNAME] [--dry-run] [--yes]`,
	Long: `Migrates the webhook configurations from version v1alpha1 to v1beta1. Version v1beta1 is the new default version
				and it is highly encouraged to use this one. Version v1alpha1 won't be supported from version Keptn 0.19.0`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return doMigration(migrateWebhooksParams)
	},
}

func doMigration(params *migrateWebhooksCmdParams) error {
	endPoint, apiToken, err := credentialmanager.NewCredentialManager(assumeYes).GetCreds(namespace)
	if err != nil {
		return fmt.Errorf(authErrorMsg)
	}

	api, err := internal.APIProvider(endPoint.String(), apiToken)
	if err != nil {
		return internal.OnAPIError(err)
	}

	if *params.ProjectName != "" { // project for migration selected
		project, err := checkProjectExists(*params.ProjectName, api)
		if err != nil {
			return fmt.Errorf("Failed to migrate webhook configurations for project %s: %s", *params.ProjectName, err.Error())
		}

		webhooks, err := getProjectWebhooks(project, api)
		if err != nil {
			return fmt.Errorf("cannot retrieve webhook-config for project %s: %s", *params.ProjectName, err.Error())
		}
		if err := migrateWebhooks(webhooks, params); err != nil {
			return fmt.Errorf("cannot migrate webhook-config for project %s: %s", *params.ProjectName, err.Error())
		}
	} else { // project for migration not selected, migrating all projects
		projects, err := api.ProjectsV1().GetAllProjects()
		if err != nil {
			return fmt.Errorf("cannot retrieve projects: %s", err.Error())
		}

		for _, project := range projects {
			webhooks, err := getProjectWebhooks(project, api)
			if err != nil {
				return fmt.Errorf("cannot retrieve webhook-config for project %s: %s", project.ProjectName, err.Error())
			}
			if err := migrateWebhooks(webhooks, params); err != nil {
				return fmt.Errorf("cannot migrate webhook-config for project %s: %s", project.ProjectName, err.Error())
			}
		}
	}

	return nil
}

func getProjectWebhooks(project *models.Project, api *api.APISet) ([]*models.Resource, error) {
	webhookURI := "%2Fwebhook%2Fwebhook.yaml"
	var webhooks []*models.Resource

	// getting project resources
	projectWebhookResource, err := api.ResourcesV1().GetProjectResource(project.ProjectName, webhookURI)
	if err != nil && err.Error() != ErrResourceNotFound.Error() {
		return nil, fmt.Errorf("cannot retrieve webhook resource on project level for project %s: %s", project.ProjectName, err.Error())
	}
	if projectWebhookResource != nil {
		webhooks = append(webhooks, projectWebhookResource)
	}

	// getting stage resources
	for _, stage := range project.Stages {
		stageWebhookResource, err := api.ResourcesV1().GetStageResource(project.ProjectName, stage.StageName, webhookURI)
		if err != nil && err.Error() != ErrResourceNotFound.Error() {
			return nil, fmt.Errorf("cannot retrieve webhook resource on stage level for project %s: %s", project.ProjectName, err.Error())
		}
		if stageWebhookResource != nil {
			webhooks = append(webhooks, stageWebhookResource)
		}
	}

	// getting service resources
	for _, stage := range project.Stages {
		for _, service := range stage.Services {
			serviceWebhookResource, err := api.ResourcesV1().GetServiceResource(project.ProjectName, stage.StageName, service.ServiceName, webhookURI)
			if err != nil && err.Error() != ErrResourceNotFound.Error() {
				return nil, fmt.Errorf("cannot retrieve webhook resource on service level for project %s: %s", project.ProjectName, err.Error())
			}
			if serviceWebhookResource != nil {
				webhooks = append(webhooks, serviceWebhookResource)
			}
		}
	}

	return webhooks, nil
}

func migrateWebhooks(webhooks []*models.Resource, params *migrateWebhooksCmdParams) error {
	fmt.Println(len(webhooks))
	// tmpWebhook := resourceToWebhook(serviceWebhookResource)
	// if tmpWebhook == nil {
	// 	return nil, fmt.Errorf("cannot decode webhook resource on service level for project %s", project.ProjectName)
	// }
	return nil
}

func resourceToWebhook(resource *models.Resource) *lib.WebHookConfig {
	whConfig, err := lib.DecodeWebHookConfigYAML([]byte(resource.ResourceContent))
	if err != nil {
		return nil
	}
	return whConfig
}

func checkProjectExists(projectName string, api *api.APISet) (*models.Project, error) {
	project, err := api.ProjectsV1().GetProject(models.Project{
		ProjectName: projectName,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve project: %s", *err.Message)
	}

	if project == nil {
		return project, fmt.Errorf("project does not exist")
	}
	return project, nil
}

func init() {
	rootCmd.AddCommand(migrateWebhooksCmd)
	migrateWebhooksParams = &migrateWebhooksCmdParams{}
	migrateWebhooksParams.DryRun = migrateWebhooksCmd.Flags().BoolP("dry-run", "", false, "Executes a dry run of webhook-config migrations without updating the files")
	migrateWebhooksParams.ProjectName = migrateWebhooksCmd.Flags().StringP("project", "", "", "The project which webhook-configs will be migrated")
	migrateWebhooksCmd.Flags().BoolVarP(&migrateWebhooksParams.AcceptAll, "yes", "y", false, "Automatically accept change of all migrations")
}
