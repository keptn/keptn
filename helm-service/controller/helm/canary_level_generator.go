package helm

import (
	corev1 "k8s.io/api/core/v1"
)

// CanaryLevelGenerator is a collection of functions which are specific for the level of the canary
type CanaryLevelGenerator interface {
	GetCanaryService(originalSvc corev1.Service, project string, stageName string) (canaryService *corev1.Service)
	IsK8sResourceDuplicated() bool
	GetNamespace(project string, stage string, generated bool) string
	DeleteCanaryRelease(project string, stage string, service string) error
}
