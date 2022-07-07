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
	t.Run("Test_LogIngestion", Test_LogIngestion)

	// Allow components to be up and running
	time.Sleep(5 * time.Minute)

	t.Run("Test_LogForwarding", Test_LogForwarding)

	t.Run("Test_ResourceServiceBasic", Test_ResourceServiceBasic)
	t.Run("Test_ManageSecrets_CreateUpdateAndDeleteSecret", Test_ManageSecrets_CreateUpdateAndDeleteSecret)
	t.Run("Test_SequenceQueue_TriggerMultiple", Test_SequenceQueue_TriggerMultiple)

	// Removed tests of webhook failing due to minishift not having connection to the outside word
	t.Run("Test_Webhook_Alpha", Test_Webhook_Alpha)
	t.Run("Test_Webhook_OverlappingSubscriptions_Beta", Test_Webhook_OverlappingSubscriptions_Beta)
	t.Run("Test_Webhook_OverlappingSubscriptions_Alpha", Test_Webhook_OverlappingSubscriptions_Alpha)
	t.Run("Test_WebhookWithDisabledFinishedEvents_Alpha", Test_WebhookWithDisabledFinishedEvents_Alpha)
	t.Run("Test_WebhookWithDisabledStartedEvents_Beta", Test_WebhookWithDisabledStartedEvents_Beta)
	t.Run("Test_WebhookWithDisabledStartedEvents_Alpha", Test_WebhookWithDisabledStartedEvents_Alpha)
	t.Run("TTest_WebhookFailInternalAddress_Beta", Test_WebhookFailInternalAddress_Beta)
	// Added a test using the API as outside address
	t.Run("Test_Webhook_Beta_API", Test_Webhook_Beta_API)

	t.Run("Test_ProvisioningURL", Test_ProvisioningURL)

	t.Run("Test_ResourceServiceGETCommitID", Test_ResourceServiceGETCommitID)
	t.Run("Test_EvaluationGitCommitID", Test_EvaluationGitCommitID)
	t.Run("Test_SSHPublicKeyAuth", Test_SSHPublicKeyAuth)

	t.Run("Test_ZeroDownTimeTriggerSequence", Test_ZeroDownTimeTriggerSequence)

	// Platform-specific Tests
}
