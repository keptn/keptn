package go_tests

import (
	"testing"
)

func Test_K3S(t *testing.T) {
	// Common Tests
	t.Run("Test_LogIngestion", Test_LogIngestion)
	t.Run("Test_LogForwarding", Test_LogForwarding)
	t.Run("Test_SelfHealing", Test_SelfHealing)
	t.Run("Test_ResourceServiceBasic", Test_ResourceServiceBasic)
	t.Run("Test_ManageSecrets_CreateUpdateAndDeleteSecret", Test_ManageSecrets_CreateUpdateAndDeleteSecret)
	t.Run("Test_SequenceQueue_TriggerMultiple", Test_SequenceQueue_TriggerMultiple)
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
	t.Run("Test_UniformRegistration_TestAPI", Test_UniformRegistration_TestAPI)
	t.Run("Test_UniformRegistration_RegistrationOfKeptnIntegration", Test_UniformRegistration_RegistrationOfKeptnIntegration)
	t.Run("Test_UniformRegistration_RegistrationOfKeptnIntegrationMultiplePods", Test_UniformRegistration_RegistrationOfKeptnIntegrationMultiplePods)
	t.Run("Test_UniformRegistration_RegistrationOfKeptnIntegrationRemoteExecPlane", Test_UniformRegistration_RegistrationOfKeptnIntegrationRemoteExecPlane)
}
