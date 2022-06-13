package go_tests

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Openshift(t *testing.T) {

	// On the minishift tests running on Github, using the rollingUpgrade strategy lead to random failures due to the
	// shipyard controller not being available after a restart.
	err := SetRecreateUpgradeStrategyForDeployment("shipyard-controller")
	require.Nil(t, err)

	// Common Tests
	t.Run("Test_LogIngestion", Test_LogIngestion)
	t.Run("Test_LogForwarding", Test_LogForwarding)
	t.Run("Test_SequenceState", Test_SequenceState)
	t.Run("Test_SequenceStateParallelStages", Test_SequenceStateParallelStages)
	t.Run("Test_SequenceStateParallelServices", Test_SequenceStateParallelServices)
	t.Run("Test_SequenceState_RetrieveMultipleSequence", Test_SequenceState_RetrieveMultipleSequence)
	t.Run("Test_SequenceState_SequenceNotFound", Test_SequenceState_SequenceNotFound)
	t.Run("Test_SequenceState_InvalidShipyard", Test_SequenceState_InvalidShipyard)
	t.Run("Test_SequenceState_CannotRetrieveShipyard", Test_SequenceState_CannotRetrieveShipyard)
	t.Run("Test_SequenceQueue", Test_SequenceQueue)
	t.Run("Test_SequenceQueue_TriggerMultiple", Test_SequenceQueue_TriggerMultiple)
	t.Run("Test_SequenceQueue_TriggerAndDeleteProject", Test_SequenceQueue_TriggerAndDeleteProject)
	t.Run("Test_SequenceLoopIntegrationTest", Test_SequenceLoopIntegrationTest)
	t.Run("Test_SelfHealing", Test_SelfHealing)
	t.Run("Test_ResourceServiceBasic", Test_ResourceServiceBasic)
	t.Run("Test_ManageSecrets_CreateUpdateAndDeleteSecret", Test_ManageSecrets_CreateUpdateAndDeleteSecret)

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

	t.Run("Test_SequenceTimeout", Test_SequenceTimeout)
	t.Run("Test_SequenceTimeoutDelayedTask", Test_SequenceTimeoutDelayedTask)
	t.Run("Test_SequenceControl_Abort", Test_SequenceControl_Abort)
	t.Run("Test_SequenceControl_AbortQueuedSequence", Test_SequenceControl_AbortQueuedSequence)
	t.Run("Test_SequenceControl_AbortPausedSequence", Test_SequenceControl_AbortPausedSequence)
	t.Run("Test_SequenceControl_AbortPausedSequenceTaskPartiallyFinished", Test_SequenceControl_AbortPausedSequenceTaskPartiallyFinished)
	t.Run("Test_SequenceControl_AbortPausedSequenceMultipleStages", Test_SequenceControl_AbortPausedSequenceMultipleStages)
	t.Run("Test_SequenceControl_PauseAndResume", Test_SequenceControl_PauseAndResume)
	t.Run("Test_SequenceControl_PauseAndResume_2", Test_SequenceControl_PauseAndResume_2)
	if res, err := CompareServiceNameWithDeploymentName("configuration-service", "resource-service"); err == nil && res {
		t.Run("Test_ResourceServiceGETCommitID", Test_ResourceServiceGETCommitID)
		t.Run("Test_EvaluationGitCommitID", Test_EvaluationGitCommitID)
		t.Run("Test_SSHPublicKeyAuth", Test_SSHPublicKeyAuth)
	}
	t.Run("Test_ZeroDownTimeTriggerSequence", Test_ZeroDownTimeTriggerSequence)
	t.Run("Test_ProvisioningURL", Test_ProvisioningURL)

	// Platform-specific Tests
}
