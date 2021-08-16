package events

import (
	"context"
	"encoding/json"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/sliceutils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/distributor/pkg/config"
	logger "github.com/sirupsen/logrus"
	"strings"
	"sync"
)

type EventMatcher struct {
	Project string
	Stage   string
	Service string
}

func NewEventMatcherFromEnv(config config.EnvConfig) *EventMatcher {
	return &EventMatcher{
		Project: config.ProjectFilter,
		Stage:   config.StageFilter,
		Service: config.ServiceFilter,
	}
}

func NewEventMatcherFromSubscription(subscription keptnmodels.EventSubscription) *EventMatcher {
	return &EventMatcher{
		Project: strings.Join(subscription.Filter.Projects, ","),
		Stage:   strings.Join(subscription.Filter.Stages, ","),
		Service: strings.Join(subscription.Filter.Services, ","),
	}
}

func (ef EventMatcher) Matches(e cloudevents.Event) bool {
	keptnBase := &v0_2_0.EventData{}
	if err := e.DataAs(keptnBase); err != nil {
		return true
	}
	if ef.Project != "" && !sliceutils.ContainsStr(strings.Split(ef.Project, ","), keptnBase.Project) ||
		ef.Stage != "" && !sliceutils.ContainsStr(strings.Split(ef.Stage, ","), keptnBase.Stage) ||
		ef.Service != "" && !sliceutils.ContainsStr(strings.Split(ef.Service, ","), keptnBase.Service) {
		return false
	}
	return true
}

type ExecutionContext struct {
	context.Context
	Wg *sync.WaitGroup
}

func NewExecutionContext(ctx context.Context, waitGroupCount int) *ExecutionContext {
	wg := new(sync.WaitGroup)
	wg.Add(waitGroupCount)
	return &ExecutionContext{ctx, wg}
}

type ceVersion struct {
	SpecVersion string `json:"specversion"`
}

func DecodeCloudEvent(data []byte) (*cloudevents.Event, error) {
	cv := &ceVersion{}
	if err := json.Unmarshal(data, cv); err != nil {
		return nil, err
	}

	event := cloudevents.NewEvent(cv.SpecVersion)

	if err := json.Unmarshal(data, &event); err != nil {
		logger.Errorf("Could not unmarshal CloudEvent: %v", err)
		return nil, err
	}

	return &event, nil
}

type CloudEventsCache struct {
	sync.RWMutex
	cache map[string][]string
}

func NewCloudEventsCache() *CloudEventsCache {
	return &CloudEventsCache{
		cache: make(map[string][]string),
	}
}

func (c *CloudEventsCache) Add(topicName, eventID string) {
	c.Lock()
	defer c.Unlock()

	eventsForTopic := c.cache[topicName]
	for _, id := range eventsForTopic {
		if id == eventID {
			return
		}
	}

	c.cache[topicName] = append(c.cache[topicName], eventID)
}

func (c *CloudEventsCache) Get(topicName string) []string {
	c.RLock()
	defer c.RUnlock()

	cp := make([]string, len(c.cache[topicName]))
	copy(cp, c.cache[topicName])
	return cp
}

// Remove a CloudEvent with specified type from the cache
func (c *CloudEventsCache) Remove(topicName, eventID string) bool {
	c.Lock()
	defer c.Unlock()

	eventsForTopic := c.cache[topicName]
	for index, id := range eventsForTopic {
		if id == eventID {
			// found
			// make sure to store the result back in c.cache[topicName]
			c.cache[topicName] = append(eventsForTopic[:index], eventsForTopic[index+1:]...)
			return true
		}
	}
	return false
}

func (c *CloudEventsCache) Contains(topicName, eventID string) bool {
	c.RLock()
	defer c.RUnlock()

	eventsForTopic := c.cache[topicName]
	for _, id := range eventsForTopic {
		if id == eventID {
			return true
		}
	}
	return false
}

func (c *CloudEventsCache) Keep(topicName string, events []*keptnmodels.KeptnContextExtendedCE) {
	c.Lock()
	defer c.Unlock()

	eventsToKeep := []string{}
	eventsForTopic := c.cache[topicName]
	for _, cacheEventID := range eventsForTopic {
		for _, e := range events {
			if cacheEventID == e.ID {
				eventsToKeep = append(eventsToKeep, e.ID)
			}
		}
	}
	c.cache[topicName] = eventsToKeep

}

func (c *CloudEventsCache) Length(topicName string) int {
	c.RLock()
	defer c.RUnlock()
	return len(c.cache[topicName])
}
