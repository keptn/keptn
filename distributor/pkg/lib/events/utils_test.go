package events

import (
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddEvent(t *testing.T) {
	cache := NewCloudEventsCache()
	cache.Add("t1", "e1")
	cache.Add("t1", "e2")
	cache.Add("t2", "e3")

	assert.True(t, cache.Contains("t1", "e1"))
	assert.True(t, cache.Contains("t1", "e2"))
	assert.False(t, cache.Contains("t1", "e3"))
	assert.True(t, cache.Contains("t2", "e3"))
}

func TestAddEventTwice(t *testing.T) {
	cache := NewCloudEventsCache()
	cache.Add("t1", "e1")
	cache.Add("t1", "e2")
	cache.Add("t1", "e2")
	assert.Equal(t, 2, cache.Length("t1"))
	assert.Equal(t, 2, len(cache.Get("t1")))
}

func TestAddRemoveEvent(t *testing.T) {
	cache := NewCloudEventsCache()
	cache.Add("t1", "e1")
	cache.Add("t1", "e2")
	cache.Add("t1", "e3")

	assert.Equal(t, 3, cache.Length("t1"))

	cache.Remove("t1", "e1")
	assert.Equal(t, 2, cache.Length("t1"))
	assert.True(t, cache.Contains("t1", "e2"))
	assert.True(t, cache.Contains("t1", "e3"))

	cache.Remove("t1", "e3")
	assert.Equal(t, 1, cache.Length("t1"))
	assert.True(t, cache.Contains("t1", "e2"))
}

func TestKeep(t *testing.T) {
	cache := NewCloudEventsCache()
	cache.Add("t1", "e1")
	cache.Add("t1", "e2")
	cache.Add("t2", "e3")
	cache.Add("t2", "e4")
	cache.Add("t2", "e5")

	cache.Keep("t1", []*models.KeptnContextExtendedCE{ce("e2")})
	cache.Keep("t2", []*models.KeptnContextExtendedCE{ce("e3"), ce("e5")})

	assert.Equal(t, 1, cache.Length("t1"))
	assert.Equal(t, 2, cache.Length("t2"))
	assert.False(t, cache.Contains("t1", "e1"))
	assert.True(t, cache.Contains("t1", "e2"))
	assert.True(t, cache.Contains("t2", "e3"))
	assert.False(t, cache.Contains("t2", "e4"))
	assert.True(t, cache.Contains("t2", "e5"))
}

func ce(id string) *models.KeptnContextExtendedCE {
	return &models.KeptnContextExtendedCE{
		ID: id,
	}
}
