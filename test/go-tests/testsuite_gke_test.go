package go_tests

import (
	"testing"
)

func Test_GKE(t *testing.T) {
	// Common Tests
	if res, err := CompareServiceNameWithDeploymentName("configuration-service", "configuration-service"); err == nil && res {
		t.Run("Test_BackupRestoreConfigService", Test_BackupRestoreConfigService)
	} else {
		t.Run("Test_BackupRestoreResourceService", Test_BackupRestoreResourceService)
	}
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
	t.Run("Test_WebhookConfigAtProjectLevel", Test_WebhookConfigAtProjectLevel)
	t.Run("Test_WebhookConfigAtStageLevel", Test_WebhookConfigAtStageLevel)
	t.Run("Test_WebhookConfigAtServiceLevel", Test_WebhookConfigAtServiceLevel)
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
		t.Run("Test_ProxyAuth", Test_ProxyAuth)
	}
	t.Run("Test_ZeroDownTimeTriggerSequence", Test_ZeroDownTimeTriggerSequence)

	// Platform-specific Tests
	t.Run("Test_ResourceService", Test_ResourceServiceBasic)
	t.Run("Test_QualityGates", Test_QualityGates)
	t.Run("Test_QualityGates_SLIWrongFinishedPayloadSend", Test_QualityGates_SLIWrongFinishedPayloadSend)
	t.Run("Test_DeliveryAssistant", Test_DeliveryAssistant)
	// TODO add resource service backup/restore test when the git credentials bug is solved
	t.Run("Test_CustomUserManagedEndpointsTest", Test_CustomUserManagedEndpointsTest)
	t.Run("Test_ContinuousDelivery", Test_ContinuousDelivery)
	t.Run("Test_GracefulShutdown", Test_GracefulShutdown)
	t.Run("Test_UniformRegistration_TestAPI", Test_UniformRegistration_TestAPI)
	t.Run("Test_UniformRegistration_RegistrationOfKeptnIntegration", Test_UniformRegistration_RegistrationOfKeptnIntegration)
	t.Run("Test_UniformRegistration_RegistrationOfKeptnIntegrationMultiplePods", Test_UniformRegistration_RegistrationOfKeptnIntegrationMultiplePods)
	t.Run("Test_UniformRegistration_RegistrationOfKeptnIntegrationRemoteExecPlane", Test_UniformRegistration_RegistrationOfKeptnIntegrationRemoteExecPlane)
}
