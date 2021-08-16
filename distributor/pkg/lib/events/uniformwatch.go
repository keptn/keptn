package events

import (
	"context"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/distributor/pkg/lib/controlplane"
	logger "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type IUniformWatch interface {
	Start(ctx context.Context)
	GetCurrentUniformSubscriptions() []models.TopicSubscription
}

type UniformWatch struct {
	ControlPlane         *controlplane.ControlPlane
	CurrentSubscriptions []models.TopicSubscription
	mtx                  sync.Mutex
}

func (sw *UniformWatch) Start(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(10 * time.Second):
			integrationData, err := sw.ControlPlane.Ping()
			if err != nil {
				logger.Errorf("Unable to send heart beat to Keptn's control plane: %s", err.Error())
				continue
			}

			sw.setCurrentSubscriptions(integrationData.Subscriptions)

			if err != nil {
				logger.Warnf("Unable to ping Keptn's control plane: %s", err.Error())
			}
		}
	}
}

func (sw *UniformWatch) setCurrentSubscriptions(subs []models.TopicSubscription) {
	sw.mtx.Lock()
	defer sw.mtx.Unlock()
	fmt.Print("Current subscriptions: ")
	fmt.Println(len(subs))
	sw.CurrentSubscriptions = subs
}

func (sw *UniformWatch) GetCurrentUniformSubscriptions() []models.TopicSubscription {
	return sw.CurrentSubscriptions
}

func NewTestUniformWatch(subscriptions []models.TopicSubscription) *TestUniformWatch {
	t := &TestUniformWatch{subscriptions}
	return t
}

type TestUniformWatch struct {
	subscriptions []models.TopicSubscription
}

func (t *TestUniformWatch) Start(ctx context.Context) {
	// no-op
}
func (t *TestUniformWatch) GetCurrentUniformSubscriptions() []models.TopicSubscription {
	return t.subscriptions
}
