package common

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLockProject(t *testing.T) {
	LockProject("my-project")
	require.NotNil(t, projectLocks["my-project"])
}

func TestUnlockProject(t *testing.T) {
	UnlockProject("my-project")
	require.NotNil(t, projectLocks["my-project"])
}
