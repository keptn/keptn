package handler_test

import (
	"github.com/ghodss/yaml"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/v0_1_4"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/remediation-service/handler"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

func newRemediation(fileName string) *v0_1_4.Remediation {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	remediation := &v0_1_4.Remediation{}
	yaml.Unmarshal(content, remediation)
	return remediation
}

func newProblemDetails(problemTitle, rootCause string) v0_2_0.ProblemDetails {
	return v0_2_0.ProblemDetails{
		ProblemTitle: problemTitle,
		RootCause:    rootCause,
	}
}

func TestGetNextAction(t *testing.T) {

	type args struct {
		remediation    *v0_1_4.Remediation
		problemDetails v0_2_0.ProblemDetails
		actionIndex    int
	}
	tests := []struct {
		name    string
		args    args
		want    *v0_2_0.ActionInfo
		wantErr bool
	}{
		{
			"determine action - by rootCause",
			args{
				newRemediation("test/remediation.yaml"),
				newProblemDetails("", "problemType1"),
				0,
			},
			&v0_2_0.ActionInfo{
				Name:        "actionName1",
				Action:      "",
				Description: "actionDescription1",
				Value:       map[string]interface{}{"foo": "bar"},
			},
			false,
		},
		{
			"determine-action - by problemTitle",
			args{
				newRemediation("test/remediation.yaml"),
				newProblemDetails("problemType1", ""),
				0,
			},
			&v0_2_0.ActionInfo{
				Name:        "actionName1",
				Action:      "",
				Description: "actionDescription1",
				Value:       map[string]interface{}{"foo": "bar"},
			},
			false,
		},
		{
			"determine action - not found",
			args{
				newRemediation("test/remediation.yaml"),
				newProblemDetails("", ""),
				0,
			},
			nil,
			true,
		},
		{
			"determine action - action index out of bound",
			args{
				newRemediation("test/remediation.yaml"),
				newProblemDetails("", ""),
				20,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handler.GetNextAction(tt.args.remediation, tt.args.problemDetails, tt.args.actionIndex)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNextAction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNextAction() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseRemediationResource(t *testing.T) {
	type args struct {
		resource *models.Resource
	}
	tests := []struct {
		name    string
		args    args
		want    *v0_1_4.Remediation
		wantErr bool
	}{
		{"", args{
			resource: newResourceFromFile("test/remediation.yaml"),
		}, newRemediation("test/remediation.yaml"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handler.ParseRemediationResource(tt.args.resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRemediationResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRemediationResource() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func newResourceFromFile(filename string) *models.Resource {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to locate resources requested by the service: %s", err.Error())
	}

	return &models.Resource{
		Metadata:        nil,
		ResourceContent: string(content),
		ResourceURI:     nil,
	}
}
