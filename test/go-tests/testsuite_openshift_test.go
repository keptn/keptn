package go_tests

import (
	"testing"
)

func Test_Openshift(t *testing.T) {
	// Common Tests
	t.Run("Test_LogIngestion", Test_LogIngestion)
	t.Run("Test_LogForwarding", Test_LogForwarding)
	t.Run("Test_SequenceState", Test_SequenceState)
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
	t.Run("Test_Webhook", Test_Webhook)
	t.Run("Test_Webhook_OverlappingSubscriptions", Test_Webhook_OverlappingSubscriptions)
	t.Run("Test_WebhookWithDisabledFinishedEvents", Test_WebhookWithDisabledFinishedEvents)
	t.Run("Test_WebhookWithDisabledFinishedEvents", Test_WebhookWithDisabledStartedEvents)
	t.Run("Test_SequenceTimeout", Test_SequenceTimeout)
	t.Run("Test_SequenceTimeoutDelayedTask", Test_SequenceTimeoutDelayedTask)
	t.Run("Test_SequenceControl_Abort", Test_SequenceControl_Abort)
	t.Run("Test_SequenceControl_AbortQueuedSequence", Test_SequenceControl_AbortQueuedSequence)
	t.Run("Test_SequenceControl_PauseAndResume", Test_SequenceControl_PauseAndResume)
	t.Run("Test_SequenceControl_PauseAndResume_2", Test_SequenceControl_PauseAndResume_2)
	if res, err := CompareServiceWithDeployment("configuration-service", "resource-service"); err == nil && res {
		t.Run("Test_ResourceServiceGETCommitID", Test_ResourceServiceGETCommitID)
		t.Run("Test_EvaluationGitCommitID", Test_EvaluationGitCommitID)
	}
	t.Run("Test_ZeroDownTimeTriggerSequence", Test_ZeroDownTimeTriggerSequence)
	t.Run("Test_SSHPublicKeyAuth", Test_SSHPublicKeyAuth)

	// Platform-specific Tests
}
