package go_tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Openshift(t *testing.T) {

	// On the minishift tests running on Github, using the rollingUpgrade strategy lead to random failures due to the
	// shipyard controller not being available after a restart.
	err := SetRecreateUpgradeStrategyForDeployment("shipyard-controller")
	require.Nil(t, err)

	// Common Tests
	// Allow components to be up and running
	time.Sleep(5 * time.Minute)
	t.Run("Test_ResourceServiceBasic", Test_ResourceServiceBasic)
	t.Run("Test_ManageSecrets_CreateUpdateAndDeleteSecret", Test_ManageSecrets_CreateUpdateAndDeleteSecret)
	t.Run("Test_SequenceQueue_TriggerMultiple", Test_SequenceQueue_TriggerMultiple)

	t.Run("Test_ResourceServiceGETCommitID", Test_ResourceServiceGETCommitID)
	t.Run("Test_EvaluationGitCommitID", Test_EvaluationGitCommitID)
	t.Run("Test_SSHPublicKeyAuth", Test_SSHPublicKeyAuth)
	t.Run("Test_ProvisioningURL", Test_ProvisioningURL)

	// Platform-specific Tests
}
