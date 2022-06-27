package http

import (
	"github.com/keptn/go-utils/pkg/api/models"
	"sync"
)

// cache is used to store key value data
type cache struct {
	sync.RWMutex
	cache map[string][]string
}

// NewCache creates a new cache
func NewCache() *cache {
	return &cache{
		cache: make(map[string][]string),
	}
}

// Add adds a new element for a given key to the cache
func (c *cache) Add(key, element string) {
	c.Lock()
	defer c.Unlock()

	eventsForTopic := c.cache[key]
	for _, id := range eventsForTopic {
		if id == element {
			return
		}
	}

	c.cache[key] = append(c.cache[key], element)
}

// Get returns all elements for a given key from the cache
func (c *cache) Get(key string) []string {
	c.RLock()
	defer c.RUnlock()

	cp := make([]string, len(c.cache[key]))
	copy(cp, c.cache[key])
	return cp
}

// Remove removes an element for a given key from the cache
func (c *cache) Remove(key, element string) bool {
	c.Lock()
	defer c.Unlock()

	eventsForTopic := c.cache[key]
	for index, id := range eventsForTopic {
		if id == element {
			// found, make sure to store the result back in c.cache[key]
			c.cache[key] = append(eventsForTopic[:index], eventsForTopic[index+1:]...)
			return true
		}
	}
	return false
}

// Contains checks whether the given element for a topic name is contained in the cache
func (c *cache) Contains(key, element string) bool {
	c.RLock()
	defer c.RUnlock()

	return c.contains(key, element)
}

// Keep deletes all elements for a topic from the cache except the ones given by events
func (c *cache) Keep(key string, elements []string) {
	c.Lock()
	defer c.Unlock()

	// keeping 0 elements, means clearing the cache
	if len(elements) == 0 {
		c.clear(key)
	}

	// convert to raw ids without duplicates
	ids := dedup(elements)

	// if none of the ids is known cached do nothing
	if !c.containsSlice(key, ids) {
		return
	}

	currentEventsForTopic := c.cache[key]
	eventsToKeep := []string{}
	for _, idOfEventToKeep := range ids {
		for _, e := range currentEventsForTopic {
			if idOfEventToKeep == e {
				eventsToKeep = append(eventsToKeep, e)
			}
		}
	}
	c.cache[key] = eventsToKeep
}

// Lenghts returns the number of cached elements for a given topic
func (c *cache) Length(key string) int {
	c.RLock()
	defer c.RUnlock()
	return len(c.cache[key])
}

func (c *cache) clear(key string) {
	c.cache[key] = []string{}
}

func (c *cache) contains(key, element string) bool {
	eventsForTopic := c.cache[key]
	for _, id := range eventsForTopic {
		if id == element {
			return true
		}
	}
	return false
}

func (c *cache) containsSlice(key string, elements []string) bool {
	contains := false
	for _, id := range elements {
		if c.contains(key, id) {
			contains = true
			break
		}
	}
	return contains
}

func dedup(elements []string) []string {
	result := make([]string, 0, len(elements))
	temp := map[string]struct{}{}
	for _, el := range elements {
		if _, ok := temp[el]; !ok {
			temp[el] = struct{}{}
			result = append(result, el)
		}
	}
	return result
}

func ToIds(events []*models.KeptnContextExtendedCE) []string {
	ids := []string{}
	for _, e := range events {
		ids = append(ids, e.ID)
	}
	return ids
}
