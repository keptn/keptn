package events

import (
	"context"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_UniformWatchReturnsRegistrationID(t *testing.T) {
	uw := NewUniformWatch(TestControlPlane{})
	uw.RegisterListener(TestListener{})

	id := uw.Start(context.TODO())
	assert.Equal(t, "a-id", id)
}

func Test_UniformWatchUpdatesListeners(t *testing.T) {

}

type TestControlPlane struct {
}

func (t TestControlPlane) Ping() (*models.Integration, error) {
	return nil, nil
}

func (t TestControlPlane) Register() (string, error) {
	return "a-id", nil
}

func (t TestControlPlane) Unregister() error {
	return nil
}

type TestListener struct {
}

func (t TestListener) UpdateSubscriptions(subscriptions []models.EventSubscription) {
	fmt.Println("update")
}
