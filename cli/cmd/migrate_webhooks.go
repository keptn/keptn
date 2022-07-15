package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/internal"
	"github.com/keptn/keptn/cli/pkg/common"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/webhook-service/lib"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type migrateWebhooksCmdParams struct {
	ProjectName *string
	DryRun      *bool
}

var migrateWebhooksParams *migrateWebhooksCmdParams
var supportedCurlMethods = [4]string{"POST", "PUT", "GET", "HEAD"}
var webhookURI = "%2Fwebhook%2Fwebhook.yaml"

const betaApiVersion = "webhookconfig.keptn.sh/v1beta1"
const alphaApiVersion = "webhookconfig.keptn.sh/v1alpha1"

type webhookResource struct {
	Project         *string
	Stage           *string
	Service         *string
	WebhookResource *models.Resource
}

var migrateWebhooksCmd = &cobra.Command{
	Use:     "migrate-webhooks",
	Short:   `Migrates the webhook configurations from version v1alpha1 to v1beta1`,
	Example: `keptn migrate-webhooks [--project=PROJECTMNAME] [--dry-run]`,
	Long: `Migrates the webhook configurations for all projects from version v1alpha1 to v1beta1. To migrate only a single project, please use "--project" flag.
				Version v1beta1 is the new default version and it is highly encouraged to use this one. Version v1alpha1 won't be supported from Keptn 0.19.0`,
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

func getProjectWebhooks(project *models.Project, api *apiutils.APISet) ([]*webhookResource, error) {
	var webhooks []*webhookResource

	// getting project resources
	projectWebhookResource, err := api.ResourcesV1().GetProjectResource(project.ProjectName, webhookURI)
	if err != nil && !errors.Is(err, apiutils.ResourceNotFoundError) {
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
		if err != nil && !errors.Is(err, apiutils.ResourceNotFoundError) {
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
			if err != nil && !errors.Is(err, apiutils.ResourceNotFoundError) {
				return nil, fmt.Errorf("cannot retrieve webhook resource on service level for project %s: %s", project.ProjectName, err.Error())
			}
			if serviceWebhookResource != nil {
				tmpWebhookResource := webhookResource{
					Project:         &project.ProjectName,
					Stage:           &stage.StageName,
					Service:         &service.ServiceName,
					WebhookResource: serviceWebhookResource,
				}
				webhooks = append(webhooks, &tmpWebhookResource)
			}
		}
	}

	return webhooks, nil
}

func migrateWebhooks(webhooks []*webhookResource, params *migrateWebhooksCmdParams, api *apiutils.APISet) error {
	for _, w := range webhooks {
		webhook := resourceToWebhook(w.WebhookResource)
		if webhook == nil {
			logrus.Errorf("cannot decode webhook resource for project %s stage %s service %s", resolveNilPointer(w.Project), resolveNilPointer(w.Stage), resolveNilPointer(w.Service))
			continue
		}
		// migrate only webhooks in v1alpha1 version
		if webhook.ApiVersion == betaApiVersion {
			continue
		}
		migratedWebhook, err := migrateAlphaWebhook(webhook)
		if err != nil {
			logrus.Errorf("cannot migrate webhook for project %s stage %s service %s", resolveNilPointer(w.Project), resolveNilPointer(w.Stage), resolveNilPointer(w.Service))
			continue
		}
		if *migrateWebhooksParams.DryRun {
			if err := printWebhook(migratedWebhook, w); err != nil {
				return err
			}
			continue
		}
		if !assumeYes {
			if err := printWebhook(migratedWebhook, w); err != nil {
				return err
			}
			userConfirmation := common.NewUserInput().AskBool("Do you want to store this migrated webhook?", &common.UserInputOptions{AssumeYes: assumeYes})
			if !userConfirmation {
				continue
			}
		}
		if err := updateWebhookResource(w, migratedWebhook, api); err != nil {
			logrus.Errorf("cannot update webhook for project %s stage %s service %s", resolveNilPointer(w.Project), resolveNilPointer(w.Stage), resolveNilPointer(w.Service))
		}
	}

	return nil
}

func printWebhook(webhook *lib.WebHookConfig, webhookResource *webhookResource) error {
	byteWebhook, err := yaml.Marshal(webhook)
	if err != nil {
		return err
	}
	fmt.Println("---------------------------------------------------------------------")
	fmt.Printf("Project:  %s\n", resolveNilPointer(webhookResource.Project))
	fmt.Printf("Stage:    %s\n", resolveNilPointer(webhookResource.Stage))
	fmt.Printf("Service:  %s\n", resolveNilPointer(webhookResource.Service))
	fmt.Println("---------------------------------------------------------------------")
	fmt.Println(string(byteWebhook))
	fmt.Println("---------------------------------------------------------------------")

	return nil
}

func resolveNilPointer(p *string) string {
	if p != nil {
		return *p
	}
	return "undefined"
}

func migrateAlphaWebhook(webhook *lib.WebHookConfig) (*lib.WebHookConfig, error) {
	migratedWebhook := lib.WebHookConfig{}
	migratedWebhook = *webhook
	migratedWebhook.ApiVersion = betaApiVersion
	for i, w := range webhook.Spec.Webhooks {
		for j, request := range w.Requests {
			betaRequest, err := migrateAlphaRequest(fmt.Sprintf("%s", request))
			if err != nil {
				return nil, err
			}
			migratedWebhook.Spec.Webhooks[i].Requests[j] = *betaRequest
		}
	}
	return &migratedWebhook, nil
}

func migrateAlphaRequest(request string) (*lib.Request, error) {
	newRequest := lib.Request{}
	curlArr, err := parseCurl(request)
	if err != nil {
		return nil, fmt.Errorf("cannot parse curl command: %s", err.Error())
	}

	//removes `curl` from array`
	curlArr = deleteItem(curlArr, 0, false)

	//extracts URL
	url, i := extractURL(curlArr)
	curlArr = deleteItem(curlArr, i, false)
	newRequest.URL = url

	//extracts Method
	method, i := extractMethod(curlArr)
	curlArr = deleteItem(curlArr, i, true)
	newRequest.Method = method

	//extracts Payload
	payload, i := extractPayload(curlArr)
	curlArr = deleteItem(curlArr, i, true)
	newRequest.Payload = payload

	//extracts Headers
	curlArr, headers := extractHeaders(curlArr)
	newRequest.Headers = headers

	//consider the rest as Options
	newRequest.Options = strings.Join(curlArr[:], " ")

	return &newRequest, nil
}

func deleteItem(arr []string, index int, previous bool) []string {
	if index < 0 {
		return arr
	}
	if previous {
		return append(arr[:index-1], arr[index+1:]...)
	}
	return append(arr[:index], arr[index+1:]...)
}

func extractURL(array []string) (string, int) {
	for i, item := range array {
		if strings.HasPrefix(item, "http") {
			return item, i
		}
	}
	return "", -1
}

func extractMethod(array []string) (string, int) {
	for i, item := range array {
		if item == "-X" || item == "--request" {
			if i < len(array)-1 {
				return array[i+1], i + 1
			}
		}
	}
	return "GET", -1
}

func extractPayload(array []string) (string, int) {
	for i, item := range array {
		if item == "-d" || item == "--data" {
			if i < len(array)-1 {
				return array[i+1], i + 1
			}
		}
	}
	return "", -1
}

func extractHeaders(array []string) ([]string, []lib.Header) {
	var headers []lib.Header
	var newArray []string
	var indexesToDelete []int
	previous := false
	for i, item := range array {
		if item == "-H" || item == "--header" {
			header := createHeader(array[i+1])
			indexesToDelete = append(indexesToDelete, i, i+1)
			headers = append(headers, header)
			previous = true
		} else {
			if previous {
				previous = false
				continue
			}
			newArray = append(newArray, item)
		}
	}

	return newArray, headers
}

func createHeader(headerStr string) lib.Header {
	arr := strings.Split(headerStr, ":")
	return lib.Header{
		Key:   strings.ReplaceAll(arr[0], " ", ""),
		Value: strings.ReplaceAll(arr[1], " ", ""),
	}
}

func updateWebhookResource(webhook *webhookResource, webhookConfig *lib.WebHookConfig, api *apiutils.APISet) error {
	resourceURI := "/webhook/webhook.yaml"
	byteWebhook, err := yaml.Marshal(webhookConfig)
	if err != nil {
		return fmt.Errorf("cannot marshal webhookConfig: %s", err.Error())
	}

	webhookResource := webhook.WebhookResource
	webhookResource.ResourceContent = string(byteWebhook)
	webhookResource.ResourceURI = &resourceURI
	webhookResources := []*models.Resource{webhookResource}

	if webhook.Service != nil && *webhook.Service != "" {
		_, err := api.ResourcesV1().UpdateServiceResources(*webhook.Project, *webhook.Stage, *webhook.Service, webhookResources)
		if err != nil {
			return fmt.Errorf("cannot update webhook resource on service level for project %s stage %s service %s: %s", *webhook.Project, *webhook.Stage, *webhook.Service, err.Error())
		}
	} else if webhook.Stage != nil && *webhook.Stage != "" {
		_, err := api.ResourcesV1().UpdateStageResources(*webhook.Project, *webhook.Stage, webhookResources)
		if err != nil {
			return fmt.Errorf("cannot update webhook resource on stage level for project %s stage %s: %s", *webhook.Project, *webhook.Stage, err.Error())
		}
	} else {
		_, err := api.ResourcesV1().UpdateProjectResources(*webhook.Project, webhookResources)
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

func checkProjectExists(projectName string, api *apiutils.APISet) (*models.Project, error) {
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

func parseCurl(command string) ([]string, error) {
	var args []string
	state := "start"
	current := ""
	quote := "\""
	escapeNext := true
	for i := 0; i < len(command); i++ {
		c := command[i]

		if state == "quotes" {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				state = "start"
			}
			continue
		}

		if escapeNext {
			current += string(c)
			escapeNext = false
			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

		if c == '"' || c == '\'' {
			state = "quotes"
			quote = string(c)
			continue
		}

		if state == "arg" {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				state = "start"
			} else {
				current += string(c)
			}
			continue
		}

		if c != ' ' && c != '\t' {
			state = "arg"
			current += string(c)
		}
	}

	if state == "quotes" {
		return []string{}, errors.New("unclosed quote in command")
	}

	if current != "" {
		args = append(args, current)
	}

	return deleteEmpty(args), nil
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func init() {
	rootCmd.AddCommand(migrateWebhooksCmd)
	migrateWebhooksParams = &migrateWebhooksCmdParams{}
	migrateWebhooksParams.DryRun = migrateWebhooksCmd.Flags().BoolP("dry-run", "", false, "Executes a dry run of webhook-config migrations without updating the files")
	migrateWebhooksParams.ProjectName = migrateWebhooksCmd.Flags().StringP("project", "", "", "The project which webhook-configs will be migrated")
}
