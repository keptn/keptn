package common

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetKeptnSpecVersion(t *testing.T) {
	specVersion := GetKeptnSpecVersion()
	assert.Equal(t, "", specVersion)

	os.Setenv(keptnSpecVersionEnvVar, "0.2.0")
	specVersion = GetKeptnSpecVersion()
	assert.Equal(t, "0.2.0", specVersion)
}
