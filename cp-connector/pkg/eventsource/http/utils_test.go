package http

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddEvent(t *testing.T) {
	cache := NewCache()
	cache.Add("t1", "e1")
	cache.Add("t1", "e2")
	cache.Add("t2", "e3")

	assert.True(t, cache.Contains("t1", "e1"))
	assert.True(t, cache.Contains("t1", "e2"))
	assert.False(t, cache.Contains("t1", "e3"))
	assert.True(t, cache.Contains("t2", "e3"))
}

func TestAddEventTwice(t *testing.T) {
	cache := NewCache()
	cache.Add("t1", "e1")
	cache.Add("t1", "e2")
	cache.Add("t1", "e2")
	assert.Equal(t, 2, cache.Length("t1"))
	assert.Equal(t, 2, len(cache.Get("t1")))
}

func TestAddRemoveEvent(t *testing.T) {
	cache := NewCache()
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

func TestKeep_NonExistingEvent(t *testing.T) {
	cache := NewCache()
	cache.Add("t1", "e1")
	cache.Add("t1", "e2")
	cache.Add("t1", "e3")

	require.Equal(t, 3, cache.Length("t1"))
	cache.Keep("t1", []string{"e0"})
	assert.Equal(t, 3, cache.Length("t1"))
}

func TestKeep_WithDuplicates(t *testing.T) {
	cache := NewCache()
	cache.Add("t1", "e1")
	cache.Add("t1", "e2")

	require.Equal(t, 2, cache.Length("t1"))
	cache.Keep("t1", []string{"e2", "e2"})
	assert.Equal(t, 1, cache.Length("t1"))
}

func TestKeep_WithEmptyEvents(t *testing.T) {
	cache := NewCache()
	cache.Add("t1", "e1")
	cache.Add("t1", "e2")

	require.Equal(t, 2, cache.Length("t1"))
	cache.Keep("t1", []string{})
	assert.Equal(t, 0, cache.Length("t1"))
}

func TestKeep(t *testing.T) {
	cache := NewCache()
	cache.Add("t1", "e1")
	cache.Add("t1", "e2")
	cache.Add("t2", "e3")
	cache.Add("t2", "e4")
	cache.Add("t2", "e5")

	cache.Keep("t1", []string{"e2"})
	cache.Keep("t2", []string{"e3", "e5"})

	assert.Equal(t, 1, cache.Length("t1"))
	assert.Equal(t, 2, cache.Length("t2"))
	assert.False(t, cache.Contains("t1", "e1"))
	assert.True(t, cache.Contains("t1", "e2"))
	assert.True(t, cache.Contains("t2", "e3"))
	assert.False(t, cache.Contains("t2", "e4"))
	assert.True(t, cache.Contains("t2", "e5"))
}
