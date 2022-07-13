package cmd

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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
var supportedCurlMethods = [4]string{"POST", "PUT", "GET", "HEAD"}

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
	newRequest := lib.Request{}
	curlArr, err := parseCurl(request)
	if err != nil {
		return nil, fmt.Errorf("cannot parse curl command: %s", err.Error())
	}

	curlArr = deleteItem(curlArr, 0, false)

	fmt.Println("zaciatok")
	fmt.Println(curlArr)

	url, i := extractURL(curlArr)
	curlArr = deleteItem(curlArr, i, false)
	newRequest.URL = url

	method, i := extractMethod(curlArr)
	curlArr = deleteItem(curlArr, i, true)
	newRequest.Method = method

	payload, i := extractPayload(curlArr)
	curlArr = deleteItem(curlArr, i, true)
	newRequest.Payload = payload

	curlArr, headers := extractHeaders(curlArr)
	newRequest.Headers = headers

	newRequest.Options = strings.Join(curlArr[:], " ")

	fmt.Println("konec")
	fmt.Println(curlArr)

	fmt.Println(newRequest)

	return nil, nil
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
		Key:   arr[0],
		Value: arr[1],
	}
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

// func parseCurl(curl string) ([]string, error) {
// 	// b, err := ioutil.ReadFile("parsecurl.js")
// 	// if err != nil {
// 	// 	fmt.Printf(err.Error())
// 	// }
// 	// ctx := v8.NewContext()
// 	// ctx.RunScript(string(b), "parsecurl.js")
// 	// ctx.RunScript(fmt.Sprintf("const result = JSON.stringify(parseCurl(%q))", curl), "main.js")
// 	// val, err := ctx.RunScript("result", "value.js")
// 	// if err != nil {
// 	// 	fmt.Printf(err.Error())
// 	// 	return
// 	// }
// 	// fmt.Printf("%s", val)

// 	arr, err := parseCommandLine(curl)
// 	if err != nil {
// 		return []string{}, err
// 	}
// 	return arr, nil
// }

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
	migrateWebhooksCmd.Flags().BoolVarP(&migrateWebhooksParams.AcceptAll, "yes", "y", false, "Automatically accept change of all migrations")
}
