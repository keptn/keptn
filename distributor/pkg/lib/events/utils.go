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
	// decode event data
	generalEventData := &v0_2_0.EventData{}
	if err := e.DataAs(generalEventData); err != nil {
		return false
	}

	if ef.Project != "" && !sliceutils.ContainsStr(strings.Split(ef.Project, ","), generalEventData.Project) ||
		ef.Stage != "" && !sliceutils.ContainsStr(strings.Split(ef.Stage, ","), generalEventData.Stage) ||
		ef.Service != "" && !sliceutils.ContainsStr(strings.Split(ef.Service, ","), generalEventData.Service) {
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

func DecodeNATSMessage(data []byte) (*cloudevents.Event, error) {
	type ceVersion struct {
		SpecVersion string `json:"specversion"`
	}

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

	return c.contains(topicName, eventID)
}

func (c *CloudEventsCache) Keep(topicName string, events []*keptnmodels.KeptnContextExtendedCE) {
	c.Lock()
	defer c.Unlock()

	// keeping 0 events, means clearing the cache
	if len(events) == 0 {
		c.clear(topicName)
	}

	// convert to raw ids without duplicates
	ids := dedup(toIDs(events))

	// if none of the ids is known cached do nothing
	if !c.containsSlice(topicName, ids) {
		return
	}

	currentEventsForTopic := c.cache[topicName]
	eventsToKeep := []string{}
	for _, idOfEventToKeep := range ids {
		for _, e := range currentEventsForTopic {
			if idOfEventToKeep == e {
				eventsToKeep = append(eventsToKeep, e)
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

func (c *CloudEventsCache) clear(topicName string) {
	c.cache[topicName] = []string{}
}

func (c *CloudEventsCache) contains(topicName, eventID string) bool {
	eventsForTopic := c.cache[topicName]
	for _, id := range eventsForTopic {
		if id == eventID {
			return true
		}
	}
	return false
}

func (c *CloudEventsCache) containsSlice(topicName string, ids []string) bool {
	contains := false
	for _, id := range ids {
		if c.contains(topicName, id) {
			contains = true
			break
		}
	}
	return contains
}

func toIDs(events []*keptnmodels.KeptnContextExtendedCE) []string {
	ids := []string{}
	for _, e := range events {
		ids = append(ids, e.ID)
	}
	return ids
}

func dedup(ids []string) []string {
	result := make([]string, 0, len(ids))
	temp := map[string]struct{}{}
	for _, item := range ids {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
