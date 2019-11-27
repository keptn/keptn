package main

import (
	"testing"

	keptnevents "github.com/keptn/go-utils/pkg/events"
	"github.com/stretchr/testify/assert"
)

func TestValidTagsDeriving(t *testing.T) {

	problemEvent := keptnevents.ProblemEventData{
		Tags:    "keptn_service:carts, keptn_stage:dev, keptn_project:sockshop",
		Project: "",
		Stage:   "",
		Service: "",
	}

	deriveFromTags(&problemEvent)

	assert.Equal(t, "sockshop", problemEvent.Project)
	assert.Equal(t, "dev", problemEvent.Stage)
	assert.Equal(t, "carts", problemEvent.Service)
}

func TestEmptyTagsDeriving(t *testing.T) {

	problemEvent := keptnevents.ProblemEventData{
		Tags:    "",
		Project: "",
		Stage:   "",
		Service: "",
	}

	deriveFromTags(&problemEvent)

	assert.Equal(t, "", problemEvent.Project)
	assert.Equal(t, "", problemEvent.Stage)
	assert.Equal(t, "", problemEvent.Service)
}
