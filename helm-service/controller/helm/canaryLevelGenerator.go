package helm

import (
	corev1 "k8s.io/api/core/v1"
)

type CanaryLevelGenerator interface {
	GetCanaryService(originalSvc corev1.Service, project string, stageName string) (canaryService *corev1.Service)
	IsK8sResourceDuplicated() bool
	GetNamespace(project string, stage string, generated bool) string
	DeleteRelease(project string, stage string, service string, generated bool, configServiceURL string) error
}
