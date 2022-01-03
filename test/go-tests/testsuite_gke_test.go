package go_tests

import (
	"testing"
)

func Test_GKE(t *testing.T) {
	// Common Tests
	t.Run("TestLogIngestion", TestLogIngestion)
	t.Run("TestLogForwarding", TestLogForwarding)
	t.Run("TestSequenceState", TestSequenceState)
	t.Run("TestSequenceState_RetrieveMultipleSequence", TestSequenceState_RetrieveMultipleSequence)
	t.Run("TestSequenceState_SequenceNotFound", TestSequenceState_SequenceNotFound)
	t.Run("TestSequenceState_InvalidShipyard", TestSequenceState_InvalidShipyard)
	t.Run("TestSequenceState_CannotRetrieveShipyard", TestSequenceState_CannotRetrieveShipyard)
	t.Run("TestSequenceLoopIntegrationTest", TestSequenceLoopIntegrationTest)
	t.Run("TestSelfHealing", TestSelfHealing)
	t.Run("TestResourceServiceBasic", TestResourceServiceBasic)
	t.Run("TestManageSecrets_CreateUpdateAndDeleteSecret", TestManageSecrets_CreateUpdateAndDeleteSecret)
	t.Run("TestWebhook", TestWebhook)
	t.Run("TestWebhook_OverlappingSubscriptions", TestWebhook_OverlappingSubscriptions)
	t.Run("TestWebhookWithDisabledFinishedEvents", TestWebhookWithDisabledFinishedEvents)
	t.Run("TestSequenceTimeout", TestSequenceTimeout)
	t.Run("TestSequenceTimeoutDelayedTask", TestSequenceTimeoutDelayedTask)
	t.Run("TestSequenceControl_Abort", TestSequenceControl_Abort)
	t.Run("TestSequenceControl_AbortQueuedSequence", TestSequenceControl_AbortQueuedSequence)
	t.Run("TestSequenceControl_PauseAndResume", TestSequenceControl_PauseAndResume)
	t.Run("TestSequenceControl_PauseAndResume_2", TestSequenceControl_PauseAndResume_2)

	// Platform-specific Tests
	t.Run("TestQualityGates", TestQualityGates)
	t.Run("TestQualityGates_BackwardsCompatibility", TestQualityGates_BackwardsCompatibility)
	t.Run("TestDeliveryAssistant", TestDeliveryAssistant)
	t.Run("TestBackupRestore", TestBackupRestore)
	t.Run("TestCustomUserManagedEndpointsTest", TestCustomUserManagedEndpointsTest)
	t.Run("TestContinuousDelivery (in-cluster/remote execution plane)", TestContinuousDelivery)
	t.Run("TestGracefulShutdown", TestGracefulShutdown)
	t.Run("TestUniformRegistration_TestAPI", TestUniformRegistration_TestAPI)
	t.Run("TestUniformRegistration_RegistrationOfKeptnIntegration", TestUniformRegistration_RegistrationOfKeptnIntegration)
	t.Run("TestUniformRegistration_RegistrationOfKeptnIntegrationMultiplePods", TestUniformRegistration_RegistrationOfKeptnIntegrationMultiplePods)
	t.Run("TestUniformRegistration_RegistrationOfKeptnIntegrationRemoteExecPlane", TestUniformRegistration_RegistrationOfKeptnIntegrationRemoteExecPlane)
}
