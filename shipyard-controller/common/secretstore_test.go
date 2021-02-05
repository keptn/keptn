package common

import (
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestCreateSecret(t *testing.T) {
	secretStore := K8sSecretStore{client: fake.NewSimpleClientset()}
	secretKey := "my-secret"
	secretVal := map[string][]byte{"git": []byte{0x1}}

	// given no secret currently stored
	secret, _ := secretStore.GetSecret(secretKey)
	assert.Empty(t, secret)

	// when creating new secret
	err := secretStore.CreateSecret(secretKey, secretVal)

	// then secret was stored
	assert.Nil(t, err)
	fetchedSecret, _ := secretStore.GetSecret(secretKey)
	assert.Equal(t, secretVal, fetchedSecret)

}

func TestCreateSecretTwice(t *testing.T) {
	secretStore := K8sSecretStore{client: fake.NewSimpleClientset()}
	secretKey := "my-secret"
	secretVal := map[string][]byte{"git": []byte{0x1}}

	// given secret is created
	secretStore.CreateSecret(secretKey, secretVal)

	// when creating same secret again
	err := secretStore.CreateSecret(secretKey, secretVal)

	// then error is returned
	assert.NotNil(t, err)
}

func TestGetSecret(t *testing.T) {
	secretStore := K8sSecretStore{client: fake.NewSimpleClientset()}
	secretKey := "my-secret"
	secretVal := map[string][]byte{"git": []byte{0x1}}

	// given secret is created
	secretStore.CreateSecret(secretKey, secretVal)

	// when getting secret
	fetchedSecret, err := secretStore.GetSecret(secretKey)

	// then stored secret is returned
	assert.Nil(t, err)
	assert.Equal(t, secretVal, fetchedSecret)
}

func TestGetNonExistentSecret(t *testing.T) {
	secretStore := K8sSecretStore{client: fake.NewSimpleClientset()}
	secretKey := "my-secret"

	// given no secret is stored
	// when getting any secret
	secret, err := secretStore.GetSecret(secretKey)

	// then no error and no secret is returned
	assert.Nil(t, secret)
	assert.Nil(t, err)
}

func TestUpdateSecret(t *testing.T) {
	secretStore := K8sSecretStore{client: fake.NewSimpleClientset()}
	secretKey := "my-secret"
	secretVal := map[string][]byte{"git": []byte{0x1}}
	updatedSecretVal := map[string][]byte{"git": []byte{0x2}}

	// given secret created
	secretStore.CreateSecret(secretKey, secretVal)
	fetchedSecret, _ := secretStore.GetSecret(secretKey)
	assert.Equal(t, secretVal, fetchedSecret)

	// when updating secret
	err := secretStore.UpdateSecret(secretKey, updatedSecretVal)

	// then secret was updated
	fetchedSecret, _ = secretStore.GetSecret(secretKey)
	assert.Nil(t, err)
	assert.Equal(t, updatedSecretVal, fetchedSecret)

}

func TestUpdateNonExistentSecret(t *testing.T) {
	secretStore := K8sSecretStore{client: fake.NewSimpleClientset()}
	secretKey := "my-secret"
	secretVal := map[string][]byte{"git": []byte{0x1}}

	// when updating secret
	err := secretStore.UpdateSecret(secretKey, secretVal)
	assert.Nil(t, err)
	fetchedSecret, _ := secretStore.GetSecret(secretKey)

	assert.Equal(t, secretVal, fetchedSecret)

}
