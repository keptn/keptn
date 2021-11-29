package common

import (
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	"testing"
)

func Test_K8sDistributedLocker(t *testing.T) {
	cm := &unstructured.Unstructured{Object: map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata": map[string]interface{}{
			"name":      "shipyard-controller-locks",
			"namespace": "keptn",
		},
	}}
	client := dynamicfake.NewSimpleDynamicClient(runtime.NewScheme(), cm)

	locker := GetK8sDistributedLockerInstance(client)

	lockID, err := locker.Lock("my-key")
	require.Nil(t, err)
	require.NotEmpty(t, lockID)

	err = locker.Unlock(lockID)
	require.Nil(t, err)
}
