package go_tests

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const provisioningShipyard = `YXBpVmVyc2lvbjogInNwZWMua2VwdG4uc2gvMC4yLjMiCmtpbmQ6ICJTaGlweWFyZCIKbWV0YWRhdGE6CiAgbmFtZTogInNoaXB5YXJkLXBvZHRhdG8tb2hlYWQiCnNwZWM6CiAgc3RhZ2VzOgogICAgLSBuYW1lOiAiZGV2IgogICAgICBzZXF1ZW5jZXM6CiAgICAgICAgLSBuYW1lOiAiZGVsaXZlcnkiCiAgICAgICAgICB0YXNrczoKICAgICAgICAgICAgLSBuYW1lOiAiZGVwbG95bWVudCIKICAgICAgICAgICAgICBwcm9wZXJ0aWVzOgogICAgICAgICAgICAgICAgZGVwbG95bWVudHN0cmF0ZWd5OiAiZGlyZWN0IgogICAgICAgICAgICAtIG5hbWU6ICJ0ZXN0IgogICAgICAgICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICAgICAgICB0ZXN0c3RyYXRlZ3k6ICJmdW5jdGlvbmFsIgogICAgICAgICAgICAtIG5hbWU6ICJldmFsdWF0aW9uIgogICAgICAgICAgICAtIG5hbWU6ICJyZWxlYXNlIgogICAgICAgIC0gbmFtZTogImRlbGl2ZXJ5LWRpcmVjdCIKICAgICAgICAgIHRhc2tzOgogICAgICAgICAgICAtIG5hbWU6ICJkZXBsb3ltZW50IgogICAgICAgICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICAgICAgICBkZXBsb3ltZW50c3RyYXRlZ3k6ICJkaXJlY3QiCiAgICAgICAgICAgIC0gbmFtZTogInJlbGVhc2UiCgogICAgLSBuYW1lOiAicHJvZCIKICAgICAgc2VxdWVuY2VzOgogICAgICAgIC0gbmFtZTogImRlbGl2ZXJ5IgogICAgICAgICAgdHJpZ2dlcmVkT246CiAgICAgICAgICAgIC0gZXZlbnQ6ICJkZXYuZGVsaXZlcnkuZmluaXNoZWQiCiAgICAgICAgICB0YXNrczoKICAgICAgICAgICAgLSBuYW1lOiAiZGVwbG95bWVudCIKICAgICAgICAgICAgICBwcm9wZXJ0aWVzOgogICAgICAgICAgICAgICAgZGVwbG95bWVudHN0cmF0ZWd5OiAiYmx1ZV9ncmVlbl9zZXJ2aWNlIgogICAgICAgICAgICAtIG5hbWU6ICJ0ZXN0IgogICAgICAgICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICAgICAgICB0ZXN0c3RyYXRlZ3k6ICJwZXJmb3JtYW5jZSIKICAgICAgICAgICAgLSBuYW1lOiAiZXZhbHVhdGlvbiIKICAgICAgICAgICAgLSBuYW1lOiAicmVsZWFzZSIKICAgICAgICAtIG5hbWU6ICJyb2xsYmFjayIKICAgICAgICAgIHRyaWdnZXJlZE9uOgogICAgICAgICAgICAtIGV2ZW50OiAicHJvZC5kZWxpdmVyeS5maW5pc2hlZCIKICAgICAgICAgICAgICBzZWxlY3RvcjoKICAgICAgICAgICAgICAgIG1hdGNoOgogICAgICAgICAgICAgICAgICByZXN1bHQ6ICJmYWlsIgogICAgICAgICAgdGFza3M6CiAgICAgICAgICAgIC0gbmFtZTogInJvbGxiYWNrIgoKICAgICAgICAtIG5hbWU6ICJkZWxpdmVyeS1kaXJlY3QiCiAgICAgICAgICB0cmlnZ2VyZWRPbjoKICAgICAgICAgICAgLSBldmVudDogImRldi5kZWxpdmVyeS1kaXJlY3QuZmluaXNoZWQiCiAgICAgICAgICB0YXNrczoKICAgICAgICAgICAgLSBuYW1lOiAiZGVwbG95bWVudCIKICAgICAgICAgICAgICBwcm9wZXJ0aWVzOgogICAgICAgICAgICAgICAgZGVwbG95bWVudHN0cmF0ZWd5OiAiZGlyZWN0IgogICAgICAgICAgICAtIG5hbWU6ICJyZWxlYXNlIg==`
const baseProjectPath = "/controlPlane/v1/project"

func Test_ProvisioningURL(t *testing.T) {
	projectName := "url-provisioning"
	keptnNamespace := GetKeptnNameSpaceFromEnv()

	t.Logf("Setting up AUTOMATIC_PROVISIONING_URL env variable")
	mockServerIP, err := GetServiceExternalIP(keptnNamespace, "mockserver")
	require.Nil(t, err)
	t.Logf("External mockserver IP address: %s", mockServerIP)
	_, err = ExecuteCommandf("kubectl set env deployment/shipyard-controller AUTOMATIC_PROVISIONING_URL=http://%s:1080 -n %s", mockServerIP, keptnNamespace)
	require.Nil(t, err)

	t.Logf("Sleeping for 30s...")
	time.Sleep(30 * time.Second)
	t.Logf("Continue to work...")

	t.Logf("Creating a new upstream repository for project %s", projectName)
	err = RecreateProjectUpstream(projectName)
	require.Nil(t, err)

	t.Logf("Creating a new project %s with a provisioned Gitea Upstream", projectName)
	projectParams := map[string]string{
		"name":     projectName,
		"shipyard": provisioningShipyard,
	}
	createProjectRequestData, err := json.Marshal(projectParams)
	require.Nil(t, err)

	resp, err := ApiPOSTRequest(baseProjectPath, createProjectRequestData, 3)
	require.Nil(t, err)
	require.Equal(t, 201, resp.Response().StatusCode)

	t.Logf("Deleting project %s with a provisioned Gitea Upstream", projectName)

	resp, err = ApiDELETERequest(baseProjectPath+"/"+projectName, 3)
	require.Nil(t, err)
	require.Equal(t, 200, resp.Response().StatusCode)

	t.Logf("Unsetting AUTOMATIC_PROVISIONING_URL env variable")
	_, err = ExecuteCommandf("kubectl set env deployment/shipyard-controller AUTOMATIC_PROVISIONING_URL=%s -n %s", "", keptnNamespace)
	require.Nil(t, err)

	t.Logf("Sleeping for 30s...")
	time.Sleep(30 * time.Second)
}
