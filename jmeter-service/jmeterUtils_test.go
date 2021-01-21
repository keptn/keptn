package main

import (
	"encoding/json"
	"fmt"
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func Test_executeJMeter(t *testing.T) {
	var returnedStatus int
	var returnedResources keptnapimodels.Resources

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			if strings.HasSuffix(r.URL.Path, "/resource") {
				marshal, _ := json.Marshal(returnedResources)
				w.Write(marshal)
				return
			}

			w.WriteHeader(returnedStatus)
			w.Write([]byte(`{
	"code": ` + fmt.Sprintf("%d", returnedStatus) + `,
	"message": ""
}`))
		}),
	)
	defer ts.Close()

	os.Setenv("CONFIGURATION_SERVICE", ts.URL)
	os.Setenv("env", "production")

	type args struct {
		testInfo       *TestInfo
		workload       *Workload
		resultsDir     string
		url            *url.URL
		LTN            string
		funcValidation bool
		logger         *keptncommon.Logger
	}
	tests := []struct {
		name              string
		args              args
		returnedStatus    int
		returnedResources []string
		want              bool
		wantErr           bool
	}{
		{
			name: "Skip tests if 404 is returned by configuration service and mark as success",
			args: args{
				testInfo: &TestInfo{
					Project:      "sockshop",
					Stage:        "dev",
					Service:      "carts",
					TestStrategy: "functional",
				},
				workload: &Workload{
					TestStrategy:      "",
					VUser:             0,
					LoopCount:         0,
					ThinkTime:         0,
					Script:            "test.jmx",
					AcceptedErrorRate: 0,
					AvgRtValidation:   0,
					Properties:        nil,
				},
				resultsDir:     "",
				url:            nil,
				LTN:            "",
				funcValidation: false,
				logger:         nil,
			},
			want:           true,
			wantErr:        false,
			returnedStatus: 404,
		},
		{
			name: "Skip tests if error code is returned by configuration service and return error",
			args: args{
				testInfo: &TestInfo{
					Project:      "sockshop",
					Stage:        "dev",
					Service:      "carts",
					TestStrategy: "functional",
				},
				workload: &Workload{
					TestStrategy:      "",
					VUser:             0,
					LoopCount:         0,
					ThinkTime:         0,
					Script:            "test.jmx",
					AcceptedErrorRate: 0,
					AvgRtValidation:   0,
					Properties:        nil,
				},
				resultsDir:     "",
				url:            nil,
				LTN:            "",
				funcValidation: false,
				logger:         nil,
			},
			want:           false,
			wantErr:        true,
			returnedStatus: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			returnedStatus = tt.returnedStatus

			returnedResources = keptnapimodels.Resources{
				Resources: []*keptnapimodels.Resource{},
			}

			for _, res := range tt.returnedResources {
				returnedResources.Resources = append(returnedResources.Resources, &keptnapimodels.Resource{
					ResourceURI: &res,
				})
			}
			got, err := executeJMeter(tt.args.testInfo, tt.args.workload, tt.args.resultsDir, tt.args.url, tt.args.LTN, tt.args.funcValidation, tt.args.logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("executeJMeter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("executeJMeter() got = %v, want %v", got, tt.want)
			}
		})
	}
}
