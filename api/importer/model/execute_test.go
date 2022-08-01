package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManifestExecution_GetProject(t *testing.T) {
	const project = "TestProjectID"
	mc := NewManifestExecution(project)
	assert.Equal(t, project, mc.GetProject())
	assert.NotNil(t, mc.Tasks)
	assert.NotNil(t, mc.Inputs)
	assert.Contains(t, mc.Inputs, projectInputContextKey)
}
