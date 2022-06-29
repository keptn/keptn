package importer

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/keptn/keptn/api/importer/model"
	"github.com/keptn/keptn/api/test/utils"
)

func TestUnmarshalManifest(t *testing.T) {
	tests := []struct {
		name                  string
		inputManifest         io.Reader
		expectedManifest      *model.ImportManifest
		expectErr             bool
		expectedError         error
		expectedErrorContains string
	}{
		{
			name:          "Basic empty manifest",
			inputManifest: strings.NewReader(`apiVersion: v1beta1`),
			expectedManifest: &model.ImportManifest{
				ApiVersion: "v1beta1",
				Tasks:      nil,
			},
			expectErr:             false,
			expectedError:         nil,
			expectedErrorContains: "",
		},
		{
			name: "Single api task manifest",
			inputManifest: strings.NewReader(
				`
                apiVersion: v1beta1
                tasks:
                  - id: sample-api
                    type: api
                    name: Sample API Task
                    payload: "api/some-payload.json"
                    action: "keptn-api-v1-endpoint-operation"
                `,
			),
			expectedManifest: &model.ImportManifest{
				ApiVersion: "v1beta1",
				Tasks: []*model.ManifestTask{
					{
						APITask: &model.APITask{
							Action:      "keptn-api-v1-endpoint-operation",
							PayloadFile: "api/some-payload.json",
						},
						ResourceTask: nil,
						ID:           "sample-api",
						Type:         "api",
						Name:         "Sample API Task",
					},
				},
			},
			expectErr:             false,
			expectedError:         nil,
			expectedErrorContains: "",
		},
		{
			name: "API and resource task manifest",
			inputManifest: strings.NewReader(
				`
                apiVersion: v1beta1
                tasks:
                  - id: sample-api
                    type: api
                    name: Sample API Task
                    payload: "api/some-payload.json"
                    action: "keptn-api-v1-endpoint-operation"
                  - id: sample-resource
                    type: resource
                    name: Sample resource task
                    resource: "resources/webhook.yaml"    # where is the file stored in the package
                    resourceUri: "webhook.yaml"           # what should the file be called in the upstream repo
                    stage: "dev"
                `,
			),
			expectedManifest: &model.ImportManifest{
				ApiVersion: "v1beta1",
				Tasks: []*model.ManifestTask{
					{
						APITask: &model.APITask{
							Action:      "keptn-api-v1-endpoint-operation",
							PayloadFile: "api/some-payload.json",
						},
						ResourceTask: nil,
						ID:           "sample-api",
						Type:         "api",
						Name:         "Sample API Task",
					},
					{
						APITask: nil,
						ResourceTask: &model.ResourceTask{
							File:      "resources/webhook.yaml",
							RemoteURI: "webhook.yaml",
							Stage:     "dev",
						},
						ID:   "sample-resource",
						Type: "resource",
						Name: "Sample resource task",
					},
				},
			},
			expectErr:             false,
			expectedError:         nil,
			expectedErrorContains: "",
		},
		{
			name: "Resource task without stage indication manifest",
			inputManifest: strings.NewReader(
				`
                apiVersion: v1beta1
                tasks:
                  - id: sample-resource
                    type: resource
                    name: Sample resource task
                    resource: "resources/webhook.yaml"    # where is the file stored in the package
                    resourceUri: "webhook.yaml"           # what should the file be called in the upstream repo
                `,
			),
			expectedManifest: &model.ImportManifest{
				ApiVersion: "v1beta1",
				Tasks: []*model.ManifestTask{
					{
						APITask: nil,
						ResourceTask: &model.ResourceTask{
							File:      "resources/webhook.yaml",
							RemoteURI: "webhook.yaml",
						},
						ID:   "sample-resource",
						Type: "resource",
						Name: "Sample resource task",
					},
				},
			},
			expectErr:             false,
			expectedError:         nil,
			expectedErrorContains: "",
		},
		{
			name:                  "Return error when reading manifest fails",
			inputManifest:         utils.NewTestReader([]byte("somerandomdatabeforeerror"), 0, true),
			expectedManifest:      nil,
			expectErr:             true,
			expectedError:         nil,
			expectedErrorContains: "error reading manifest: ",
		},
		{
			name:                  "Return error when manifest is not valid YAML",
			inputManifest:         strings.NewReader("some random string without any yaml resemblance"),
			expectedManifest:      nil,
			expectErr:             true,
			expectedError:         nil,
			expectedErrorContains: "error unmarshalling yaml manifest: ",
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				sut := new(YAMLManifestUnMarshaler)
				manifest, err := sut.Parse(test.inputManifest)
				assert.Equal(t, test.expectedManifest, manifest)
				if test.expectErr {
					assert.Error(t, err)
				}
				if test.expectedError != nil {
					assert.ErrorIs(t, err, test.expectedError)
				}
				if test.expectedErrorContains != "" {
					assert.ErrorContains(t, err, test.expectedErrorContains)
				}
			},
		)
	}
}
