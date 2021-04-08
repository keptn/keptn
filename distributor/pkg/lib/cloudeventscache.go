package lib

import (
	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
	"sync"
)

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
