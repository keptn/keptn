package go_tests

import (
	"context"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func Test_ManageSecrets_CreateUpdateAndDeleteSecret(t *testing.T) {
	k8s := KubeClient(t)
	var ns = GetKeptnNameSpaceFromEnv()
	secret1 := "my-new-secret"
	secret2 := "my-new-secret-2"

	// create secret 1
	_, err := ExecuteCommandf("keptn create secret %s --from-literal=mykey1=myvalue1", secret1)
	require.Nil(t, err)

	// create secret 2
	_, err = ExecuteCommandf("keptn create secret %s --from-literal=mykey2=myvalue2", secret2)
	require.Nil(t, err)

	// check created k8s secret 1
	k8sSecret1, err := k8s.CoreV1().Secrets(ns).Get(context.TODO(), secret1, v1.GetOptions{})
	require.Nil(t, err)
	require.Equal(t, "myvalue1", string(k8sSecret1.Data["mykey1"]))

	// check created k8s secret 2
	k8sSecret2, err := k8s.CoreV1().Secrets(ns).Get(context.TODO(), secret2, v1.GetOptions{})
	require.Nil(t, err)
	require.Equal(t, "myvalue2", string(k8sSecret2.Data["mykey2"]))

	// update secret 1
	_, err = ExecuteCommandf("keptn update secret %s --from-literal=mykey1=changed-value", secret1)
	require.Nil(t, err)

	// check update of k8s secret 1
	k8sSecret1, err = k8s.CoreV1().Secrets(ns).Get(context.TODO(), secret1, v1.GetOptions{})
	require.Nil(t, err)
	require.Equal(t, "changed-value", string(k8sSecret1.Data["mykey1"]))

	// check created k8s roles
	role, _ := k8s.RbacV1().Roles(ns).Get(context.TODO(), "keptn-secrets-default-read", v1.GetOptions{})
	require.Contains(t, role.Rules[0].ResourceNames, secret1)
	require.Contains(t, role.Rules[0].ResourceNames, secret2)

	_, err = ExecuteCommandf("keptn delete secret %s", secret1)
	require.Nil(t, err)

	_, err = ExecuteCommandf("keptn delete secret %s", secret2)
	require.Nil(t, err)

	role, _ = k8s.RbacV1().Roles(ns).Get(context.TODO(), "keptn-secrets-default-read", v1.GetOptions{})
	require.NotContains(t, role.Rules[0].ResourceNames, secret1)
}
