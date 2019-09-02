package helm

import (
	keptnevents "github.com/keptn/go-utils/pkg/events"
	corev1 "k8s.io/api/core/v1"
)

type CanaryLevelGenerator interface {
	GetCanaryService(originalSvc corev1.Service, event *keptnevents.ServiceCreateEventData, stageName string) (canaryService *corev1.Service)
	IsK8sResourceDuplicated() bool
	GetNamespace(project string, stage string, generated bool) string
	DeleteRelease(project string, service string, stage string, generated bool, configServiceURL string) error
}
