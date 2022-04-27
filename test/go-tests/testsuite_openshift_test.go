package go_tests

import (
	"context"
	keptnkubeutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

func Test_Openshift(t *testing.T) {

	clientset, err := keptnkubeutils.GetClientset(false)
	require.Nil(t, err)

	shipyardDeployment, err := clientset.AppsV1().Deployments(GetKeptnNameSpaceFromEnv()).Get(context.TODO(), "shipyard-controller", v1.GetOptions{})
	require.Nil(t, err)

	shipyardDeployment.Spec.Strategy.Type = "Recreate"
	shipyardDeployment.Spec.Strategy.RollingUpdate = nil

	_, err = clientset.AppsV1().Deployments(GetKeptnNameSpaceFromEnv()).Update(context.TODO(), shipyardDeployment, v1.UpdateOptions{})
	require.Nil(t, err)

	<-time.After(1 * time.Minute)

	// Common Tests
	t.Run("Test_LogIngestion", Test_LogIngestion)
	t.Run("Test_LogForwarding", Test_LogForwarding)
	t.Run("Test_SequenceState", Test_SequenceState)
	t.Run("Test_SequenceStateParallelStages", Test_SequenceStateParallelStages)
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
	t.Run("Test_Webhook_Alpha", Test_Webhook_Alpha)
	t.Run("Test_Webhook_Beta", Test_Webhook_Beta)
	t.Run("Test_Webhook_OverlappingSubscriptions_Alpha", Test_Webhook_OverlappingSubscriptions_Alpha)
	t.Run("Test_Webhook_OverlappingSubscriptions_Beta", Test_Webhook_OverlappingSubscriptions_Beta)
	t.Run("Test_WebhookWithDisabledFinishedEvents_Alpha", Test_WebhookWithDisabledFinishedEvents_Alpha)
	t.Run("Test_WebhookWithDisabledFinishedEvents_Beta", Test_WebhookWithDisabledFinishedEvents_Beta)
	t.Run("Test_WebhookWithDisabledStartedEvents_Alpha", Test_WebhookWithDisabledStartedEvents_Alpha)
	t.Run("Test_WebhookWithDisabledStartedEvents_Beta", Test_WebhookWithDisabledStartedEvents_Beta)
	t.Run("Test_WebhookConfigAtProjectLevel_Alpha", Test_WebhookConfigAtProjectLevel_Alpha)
	t.Run("Test_WebhookConfigAtProjectLevel_Beta", Test_WebhookConfigAtProjectLevel_Beta)
	t.Run("Test_WebhookConfigAtStageLevel_Alpha", Test_WebhookConfigAtStageLevel_Alpha)
	t.Run("Test_WebhookConfigAtStageLevel_Beta", Test_WebhookConfigAtStageLevel_Beta)
	t.Run("Test_WebhookConfigAtServiceLevel_Alpha", Test_WebhookConfigAtServiceLevel_Alpha)
	t.Run("Test_WebhookConfigAtServiceLevel_Beta", Test_WebhookConfigAtServiceLevel_Beta)
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
	t.Run("Test_ProvisioningURL", Test_ProvisioningURL)

	// Platform-specific Tests
}
