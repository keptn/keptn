package go_tests

import (
	"testing"
)

func Test_GKE(t *testing.T) {
	// Shut DownTests
	if res, err := CompareServiceNameWithDeploymentName("configuration-service", "configuration-service"); err == nil && res {
		t.Run("Test_BackupRestoreConfigService", Test_BackupRestoreConfigService)
	} else {
		t.Run("Test_BackupRestoreResourceService", Test_BackupRestoreResourceService)
	}
	t.Run("Test_GracefulShutdown", Test_GracefulShutdown)
	// Common Tests
	t.Run("Test_LogIngestion", Test_LogIngestion)
	t.Run("Test_LogForwarding", Test_LogForwarding)
	t.Run("Test_SelfHealing", Test_SelfHealing)
	t.Run("Test_ResourceServiceBasic", Test_ResourceServiceBasic)
	t.Run("Test_ManageSecrets_CreateUpdateAndDeleteSecret", Test_ManageSecrets_CreateUpdateAndDeleteSecret)
	t.Run("Test_SequenceQueue_TriggerMultiple", Test_SequenceQueue_TriggerMultiple)
	t.Run("Test_Webhook_Beta", Test_Webhook_Beta)
	t.Run("Test_Webhook_Alpha", Test_Webhook_Alpha)
	t.Run("Test_Webhook_OverlappingSubscriptions_Beta", Test_Webhook_OverlappingSubscriptions_Beta)
	t.Run("Test_Webhook_OverlappingSubscriptions_Alpha", Test_Webhook_OverlappingSubscriptions_Alpha)
	t.Run("Test_WebhookWithDisabledFinishedEvents_Beta", Test_WebhookWithDisabledFinishedEvents_Beta)
	t.Run("Test_WebhookWithDisabledFinishedEvents_Alpha", Test_WebhookWithDisabledFinishedEvents_Alpha)
	t.Run("Test_WebhookWithDisabledStartedEvents_Beta", Test_WebhookWithDisabledStartedEvents_Beta)
	t.Run("Test_WebhookWithDisabledStartedEvents_Alpha", Test_WebhookWithDisabledStartedEvents_Alpha)
	t.Run("Test_WebhookConfigAtProjectLevel_Beta", Test_WebhookConfigAtProjectLevel_Beta)
	t.Run("Test_WebhookConfigAtProjectLevel_Alpha", Test_WebhookConfigAtProjectLevel_Alpha)
	t.Run("Test_WebhookConfigAtStageLevel_Beta", Test_WebhookConfigAtStageLevel_Beta)
	t.Run("Test_WebhookConfigAtStageLevel_Alpha", Test_WebhookConfigAtStageLevel_Alpha)
	t.Run("Test_WebhookConfigAtServiceLevel_Beta", Test_WebhookConfigAtServiceLevel_Beta)
	t.Run("Test_WebhookConfigAtServiceLevel_Alpha", Test_WebhookConfigAtServiceLevel_Alpha)
	t.Run("TTest_WebhookFailInternalAddress_Beta", Test_WebhookFailInternalAddress_Beta)
	t.Run("Test_ProvisioningURL", Test_ProvisioningURL)
	if res, err := CompareServiceNameWithDeploymentName("configuration-service", "resource-service"); err == nil && res {
		t.Run("Test_ResourceServiceGETCommitID", Test_ResourceServiceGETCommitID)
		t.Run("Test_EvaluationGitCommitID", Test_EvaluationGitCommitID)
		t.Run("Test_SSHPublicKeyAuth", Test_SSHPublicKeyAuth)
		t.Run("Test_ProxyAuth", Test_ProxyAuth)
	}
	t.Run("Test_ZeroDownTimeTriggerSequence", Test_ZeroDownTimeTriggerSequence)

	// Platform-specific Tests
	t.Run("Test_QualityGates", Test_QualityGates)
	t.Run("Test_QualityGates_SLIWrongFinishedPayloadSend", Test_QualityGates_SLIWrongFinishedPayloadSend)
	t.Run("Test_QualityGates_AbortedFinishedPayloadSend", Test_QualityGates_AbortedFinishedPayloadSend)
	t.Run("Test_QualityGates_ErroredFinishedPayloadSend", Test_QualityGates_ErroredFinishedPayloadSend)
	t.Run("Test_DeliveryAssistant", Test_DeliveryAssistant)
	t.Run("Test_CustomUserManagedEndpointsTest", Test_CustomUserManagedEndpointsTest)
	t.Run("Test_ContinuousDelivery", Test_ContinuousDelivery)
	t.Run("Test_UniformRegistration_TestAPI", Test_UniformRegistration_TestAPI)
	t.Run("Test_UniformRegistration_RegistrationOfKeptnIntegration", Test_UniformRegistration_RegistrationOfKeptnIntegration)
	t.Run("Test_UniformRegistration_RegistrationOfKeptnIntegrationMultiplePods", Test_UniformRegistration_RegistrationOfKeptnIntegrationMultiplePods)
	t.Run("Test_UniformRegistration_RegistrationOfKeptnIntegrationRemoteExecPlane", Test_UniformRegistration_RegistrationOfKeptnIntegrationRemoteExecPlane)
}
