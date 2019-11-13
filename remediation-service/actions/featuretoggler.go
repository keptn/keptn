package actions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnmodels "github.com/keptn/go-utils/pkg/models"
)

// FeatureToggler ...
type FeatureToggler struct {
}

// NewFeatureToggler creates a new Aborter
func NewFeatureToggler() *FeatureToggler {
	return &FeatureToggler{}
}

// GetAction return name of action
func (f FeatureToggler) GetAction() string {
	return "featuretoggle"
}

func (f FeatureToggler) ExecuteAction(problem *keptnevents.ProblemEventData, shkeptncontext string,
	action *keptnmodels.RemediationAction) error {

	if !strings.Contains(action.Value, ":") {
		return errors.New("feature toggle remediation action not well formed")
	}

	togglename := strings.Split(action.Value, ":")[0]
	togglevalue := strings.Split(action.Value, ":")[1]

	err := sendDTProblemComment(problem.PID, "Keptn triggering change of feature toggle "+togglename+" to be set to value: "+togglevalue)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = f.ToggleFeature(togglename, togglevalue)
	if err != nil {
		sendDTProblemComment(problem.PID, "Keptn could not change feature toggle "+togglename+" to be set to value: "+togglevalue)
		return err
	}

	err = sendDTProblemComment(problem.PID, "Keptn finished change of feature toggle "+togglename+" to be set to value: "+togglevalue)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (f FeatureToggler) ResolveAction(problem *keptnevents.ProblemEventData, shkeptncontext string,
	action *keptnmodels.RemediationAction) error {
	return errors.New("no resolving action for action " + f.GetAction() + " implemented")
}

func (f FeatureToggler) ToggleFeature(togglename string, togglevalue string) error {

	unleashAPIUrl := os.Getenv("UNLEASH_SERVER_URL")
	unleashUser := os.Getenv("UNLEASH_USER")
	unleashToken := os.Getenv("UNLEASH_TOKEN")
	unleashAPIUrlExt := "/admin/features/" + togglename + "/toggle/" + togglevalue

	client := &http.Client{}
	req, err := http.NewRequest("POST", unleashAPIUrl+unleashAPIUrlExt, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(unleashUser, unleashToken)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	fmt.Println("unleash status code: " + strconv.Itoa(resp.StatusCode))

	if resp.StatusCode != 200 || resp.StatusCode != 201 {
		return errors.New("could not update feature toggle")
	}

	return nil
}

func sendDTProblemComment(problemID string, comment string) error {
	dtTenant := os.Getenv("DT_TENANT")
	dtAPIToken := os.Getenv("DT_API_TOKEN")
	dtAPIUrl := "https://" + dtTenant + "/api/v1/problem/details/" + problemID + "/comments"

	dtCommentPayload := map[string]string{"comment": comment, "user": "keptn", "context": "keptn-remediation"}
	jsonPayload, _ := json.Marshal(dtCommentPayload)

	client := &http.Client{}
	req, err := http.NewRequest("POST", dtAPIUrl, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Api-Token "+dtAPIToken)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	fmt.Println("dynatrace status code: " + strconv.Itoa(resp.StatusCode))

	return nil

}
