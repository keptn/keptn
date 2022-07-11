package go_tests

import (
	"testing"
)

func Test_GKE(t *testing.T) {
	// Shut DownTests

	t.Run("Test_BackupRestore", Test_BackupRestore)
	t.Run("Test_GracefulShutdown", Test_GracefulShutdown)
	// Common Tests
	t.Run("Test_ResourceServiceBasic", Test_ResourceServiceBasic)
	t.Run("Test_ManageSecrets_CreateUpdateAndDeleteSecret", Test_ManageSecrets_CreateUpdateAndDeleteSecret)
	t.Run("Test_SequenceQueue_TriggerMultiple", Test_SequenceQueue_TriggerMultiple)
	t.Run("Test_Webhook_Failures", Test_Webhook_Failures)
	t.Run("Test_Webhook", Test_Webhook)
	t.Run("Test_ExecutingWebhookTargetingClusterInternalAddressesFails", Test_ExecutingWebhookTargetingClusterInternalAddressesFails)

	t.Run("Test_ProvisioningURL", Test_ProvisioningURL)

	t.Run("Test_ResourceServiceGETCommitID", Test_ResourceServiceGETCommitID)
	t.Run("Test_EvaluationGitCommitID", Test_EvaluationGitCommitID)
	t.Run("Test_SSHPublicKeyAuth", Test_SSHPublicKeyAuth)
	t.Run("Test_ProxyAuth", Test_ProxyAuth)

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
