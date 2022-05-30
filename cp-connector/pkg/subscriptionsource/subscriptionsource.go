package subscriptionsource

import (
	"context"
	"github.com/keptn/keptn/cp-connector/pkg/types"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cp-connector/pkg/logger"
)

type SubscriptionSource interface {
	Start(context.Context, types.RegistrationData, chan []models.EventSubscription) error
	Register(integration models.Integration) (string, error)
}

var _ SubscriptionSource = FixedSubscriptionSource{}
var _ SubscriptionSource = (*UniformSubscriptionSource)(nil)

// UniformSubscriptionSource represents a source for uniform subscriptions
type UniformSubscriptionSource struct {
	uniformAPI    api.UniformV1Interface
	clock         clock.Clock
	fetchInterval time.Duration
	logger        logger.Logger
}

func (s *UniformSubscriptionSource) Register(integration models.Integration) (string, error) {
	integrationID, err := s.uniformAPI.RegisterIntegration(integration)
	if err != nil {
		return "", err
	}
	return integrationID, nil
}

// WithFetchInterval specifies the interval the subscription source should
// use when polling for new subscriptions
func WithFetchInterval(interval time.Duration) func(s *UniformSubscriptionSource) {
	return func(s *UniformSubscriptionSource) {
		s.fetchInterval = interval
	}
}

// New creates a new UniformSubscriptionSource
func New(uniformAPI api.UniformV1Interface, options ...func(source *UniformSubscriptionSource)) *UniformSubscriptionSource {
	subscriptionSource := &UniformSubscriptionSource{uniformAPI: uniformAPI, clock: clock.New(), fetchInterval: time.Second * 5, logger: logger.NewDefaultLogger()}
	for _, o := range options {
		o(subscriptionSource)
	}
	return subscriptionSource
}

// Start triggers the execution of the UniformSubscriptionSource
func (s *UniformSubscriptionSource) Start(ctx context.Context, registrationData types.RegistrationData, subscriptionChannel chan []models.EventSubscription) error {
	ticker := s.clock.Ticker(s.fetchInterval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				updatedIntegrationData, err := s.uniformAPI.Ping(registrationData.ID)
				if err != nil {
					s.logger.Errorf("Unable to ping control plane: %v", err)
					continue
				}
				subscriptionChannel <- updatedIntegrationData.Subscriptions
			}
		}
	}()
	return nil
}

// FixedSubscriptionSource can be used to use a fixed list of subscriptions rather than
// consulting the Keptn API for subscriptions.
// This is useful when you want to consume events from an event source, but NOT register
// as an Keptn integration to the control plane
type FixedSubscriptionSource struct {
	fixedSubscriptions []models.EventSubscription
}

// WithFixedSubscriptions adds a fixed list of subscriptions to the FixedSubscriptionSource
func WithFixedSubscriptions(subscriptions ...models.EventSubscription) func(s *FixedSubscriptionSource) {
	return func(s *FixedSubscriptionSource) {
		s.fixedSubscriptions = subscriptions
	}
}

// NewFixedSubscriptionSource creates a new instance of FixedSubscriptionSource
func NewFixedSubscriptionSource(options ...func(source *FixedSubscriptionSource)) *FixedSubscriptionSource {
	fss := &FixedSubscriptionSource{fixedSubscriptions: []models.EventSubscription{}}
	for _, o := range options {
		o(fss)
	}
	return fss
}

func (s FixedSubscriptionSource) Start(ctx context.Context, data types.RegistrationData, c chan []models.EventSubscription) error {
	go func() { c <- s.fixedSubscriptions }()
	return nil
}

func (s FixedSubscriptionSource) Register(integration models.Integration) (string, error) {
	return "", nil
}
