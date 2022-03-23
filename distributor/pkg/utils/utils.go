package utils

import (
	"context"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"sync"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/sliceutils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/distributor/pkg/config"
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

func NewEventMatcherFromSubscription(subscription apimodels.EventSubscription) *EventMatcher {
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
	Wg       *sync.WaitGroup
	CancelFn context.CancelFunc
}

func NewExecutionContext(ctx context.Context, waitGroupCount int) *ExecutionContext {
	wg := new(sync.WaitGroup)
	wg.Add(waitGroupCount)
	return &ExecutionContext{ctx, wg, func() {}}
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
		return nil, err
	}

	return &event, nil
}

// Cache is used to store key value data
type Cache struct {
	sync.RWMutex
	cache map[string][]string
}

// NewCache creates a new cache
func NewCache() *Cache {
	return &Cache{
		cache: make(map[string][]string),
	}
}

// Add adds a new element for a given key to the cache
func (c *Cache) Add(key, element string) {
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
func (c *Cache) Get(key string) []string {
	c.RLock()
	defer c.RUnlock()

	cp := make([]string, len(c.cache[key]))
	copy(cp, c.cache[key])
	return cp
}

// Remove removes an element for a given key from the cache
func (c *Cache) Remove(key, element string) bool {
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
func (c *Cache) Contains(key, element string) bool {
	c.RLock()
	defer c.RUnlock()

	return c.contains(key, element)
}

// Keep deletes all elements for a topic from the cache except the ones given by events
func (c *Cache) Keep(key string, elements []string) {
	c.Lock()
	defer c.Unlock()

	// keeping 0 elements, means clearing the cache
	if len(elements) == 0 {
		c.clear(key)
	}

	// convert to raw ids without duplicates
	ids := Dedup(elements)

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
func (c *Cache) Length(key string) int {
	c.RLock()
	defer c.RUnlock()
	return len(c.cache[key])
}

func (c *Cache) clear(key string) {
	c.cache[key] = []string{}
}

func (c *Cache) contains(key, element string) bool {
	eventsForTopic := c.cache[key]
	for _, id := range eventsForTopic {
		if id == element {
			return true
		}
	}
	return false
}

func (c *Cache) containsSlice(key string, elements []string) bool {
	contains := false
	for _, id := range elements {
		if c.contains(key, id) {
			contains = true
			break
		}
	}
	return contains
}

// Dedup removes duplicate elements from the given list of strings
func Dedup(elements []string) []string {
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

func ToIds(events []*apimodels.KeptnContextExtendedCE) []string {
	ids := []string{}
	for _, e := range events {
		ids = append(ids, e.ID)
	}
	return ids
}

func IsEqual(a1 []string, a2 []string) bool {
	sort.Strings(a2)
	sort.Strings(a1)
	return reflect.DeepEqual(a1, a2)
}
