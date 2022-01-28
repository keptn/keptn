package auth

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestState(t *testing.T) {
	t.Run("State - len 10", func(t *testing.T) {
		state, err := State(10)
		require.Nil(t, err)
		decoded, err := base64.StdEncoding.DecodeString(state)
		require.Nil(t, err)
		assert.Equal(t, 10, len(decoded))
	})
	t.Run("State - len 1", func(t *testing.T) {
		state, err := State(1)
		require.Nil(t, err)
		decoded, err := base64.StdEncoding.DecodeString(state)
		require.Nil(t, err)
		assert.Equal(t, 1, len(decoded))
	})
	t.Run("State - len 0", func(t *testing.T) {
		state, err := State(0)
		require.Equal(t, "", state)
		require.NotNil(t, err)
	})
	t.Run("State - len -1", func(t *testing.T) {
		state, err := State(-1)
		require.Equal(t, "", state)
		require.NotNil(t, err)
	})

}
