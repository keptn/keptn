package go_tests

import (
	"fmt"
	"github.com/keptn/go-utils/pkg/common/osutils"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

const airgappedRegistryUrlEnvVarName = "AIRGAPPED_REGISTRY_URL"

func Test_AirgappedImagesAreSetCorrectly(t *testing.T) {
	airgappedRegistryUrl := osutils.GetOSEnv(airgappedRegistryUrlEnvVarName)
	require.NotEmpty(t, airgappedRegistryUrl)

	out, err := ExecuteCommand(fmt.Sprintf("kubectl get pods -n %s -o jsonpath=\"{.items[*].spec.containers[*].image}\"", GetKeptnNameSpaceFromEnv()))
	require.Nil(t, err)
	keptnImages := strings.Split(out, " ")

	for _, image := range keptnImages {
		if strings.HasPrefix(image, "rancher/") {
			// Built-in k3s images don't need to be checked
			continue
		}

		require.Contains(t, image, airgappedRegistryUrl)
	}
}
