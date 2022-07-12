package cmd

import (
	"encoding/base64"
	"encoding/json"
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

const webhookURI = "%2Fwebhook%2Fwebhook.yaml"
const betaApiVersion = "webhookconfig.keptn.sh/v1beta1"
const alphaApiVersion = "webhookconfig.keptn.sh/v1alpha1"

type webhookResource struct {
	Project         *string
	Stage           *string
	Sevice          *string
	WebhookResource *models.Resource
}

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
		if err := migrateWebhooks(webhooks, params, api); err != nil {
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
			if err := migrateWebhooks(webhooks, params, api); err != nil {
				return fmt.Errorf("cannot migrate webhook-config for project %s: %s", project.ProjectName, err.Error())
			}
		}
	}

	return nil
}

func getProjectWebhooks(project *models.Project, api *api.APISet) ([]*webhookResource, error) {
	var webhooks []*webhookResource

	// getting project resources
	projectWebhookResource, err := api.ResourcesV1().GetProjectResource(project.ProjectName, webhookURI)
	if err != nil && err.Error() != ErrResourceNotFound.Error() {
		return nil, fmt.Errorf("cannot retrieve webhook resource on project level for project %s: %s", project.ProjectName, err.Error())
	}
	if projectWebhookResource != nil {
		tmpWebhookResource := webhookResource{
			Project:         &project.ProjectName,
			WebhookResource: projectWebhookResource,
		}
		webhooks = append(webhooks, &tmpWebhookResource)
	}

	// getting stage resources
	for _, stage := range project.Stages {
		stageWebhookResource, err := api.ResourcesV1().GetStageResource(project.ProjectName, stage.StageName, webhookURI)
		if err != nil && err.Error() != ErrResourceNotFound.Error() {
			return nil, fmt.Errorf("cannot retrieve webhook resource on stage level for project %s: %s", project.ProjectName, err.Error())
		}
		if stageWebhookResource != nil {
			tmpWebhookResource := webhookResource{
				Project:         &project.ProjectName,
				Stage:           &stage.StageName,
				WebhookResource: stageWebhookResource,
			}
			webhooks = append(webhooks, &tmpWebhookResource)
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
				tmpWebhookResource := webhookResource{
					Project:         &project.ProjectName,
					Stage:           &stage.StageName,
					Sevice:          &service.ServiceName,
					WebhookResource: serviceWebhookResource,
				}
				webhooks = append(webhooks, &tmpWebhookResource)
			}
		}
	}

	return webhooks, nil
}

func migrateWebhooks(webhooks []*webhookResource, params *migrateWebhooksCmdParams, api *api.APISet) error {
	for _, w := range webhooks {
		webhook := resourceToWebhook(w.WebhookResource)
		if webhook == nil {
			return fmt.Errorf("cannot decode webhook resource")
		}
		// migrate only webhooks in v1alpha1 version
		if webhook.ApiVersion == betaApiVersion {
			continue
		}
		migratedWebhook, err := migrateAlphaWebhook(webhook)
		if err != nil {
			return err
		}
		if *migrateWebhooksParams.DryRun {
			byteWebhook, err := json.Marshal(migratedWebhook)
			if err != nil {
				return err
			}
			fmt.Println(string(byteWebhook))
		}
		if err := updateWebhookResources(w, migratedWebhook, api); err != nil {
			return err
		}
	}

	return nil
}

func migrateAlphaWebhook(webhook *lib.WebHookConfig) (*lib.WebHookConfig, error) {
	migratedWebhook := webhook
	migratedWebhook.ApiVersion = betaApiVersion
	for i, w := range webhook.Spec.Webhooks {
		for j, request := range w.Requests {
			betaRequest, err := migrateAlphaRequest(fmt.Sprintf("%s", request))
			if err != nil {
				return nil, err
			}
			migratedWebhook.Spec.Webhooks[i].Requests[j] = betaRequest
		}
	}
	return migratedWebhook, nil
}

func migrateAlphaRequest(request string) (*lib.Request, error) {
	return nil, nil
}

func updateWebhookResources(webhook *webhookResource, webhookConfig *lib.WebHookConfig, api *api.APISet) error {
	byteWebhook, err := json.Marshal(webhookConfig)
	if err != nil {
		return fmt.Errorf("cannot marshal webhookConfig: %s", err.Error())
	}
	encodedWebhook := base64.StdEncoding.EncodeToString(byteWebhook)
	webhookResource := webhook.WebhookResource
	webhookResource.ResourceContent = encodedWebhook

	if webhook.Sevice != nil && *webhook.Sevice != "" {
		_, err := api.ResourcesV1().UpdateServiceResource(*webhook.Project, *webhook.Stage, *webhook.Sevice, webhookResource)
		if err != nil {
			return fmt.Errorf("cannot update webhook resource on service level for project %s stage %s service %s: %s", *webhook.Project, *webhook.Stage, *webhook.Sevice, err.Error())
		}
	} else if webhook.Stage != nil && *webhook.Stage != "" {
		_, err := api.ResourcesV1().UpdateStageResource(*webhook.Project, *webhook.Stage, webhookResource)
		if err != nil {
			return fmt.Errorf("cannot update webhook resource on stage level for project %s stage %s: %s", *webhook.Project, *webhook.Stage, err.Error())
		}
	} else {
		_, err := api.ResourcesV1().UpdateProjectResource(*webhook.Project, webhookResource)
		if err != nil {
			return fmt.Errorf("cannot update webhook resource on project level for project %s: %s", *webhook.Project, err.Error())
		}
	}
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
