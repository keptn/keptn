package main

import (
	"os"
	"testing"

	"github.com/keptn/go-utils/pkg/utils"
	"github.com/magiconair/properties/assert"
)

func TestGetNextStage(t *testing.T) {

	const org = "keptn"
	const project = "examples"
	err := os.RemoveAll(project)
	assert.Equal(t, err, nil, "Received unexpected error")

	repo, err := utils.Checkout(org, project, "master")
	assert.Equal(t, err, nil, "Received unexpected error")

	data := evaluationDoneEvent{Project: "examples/onboarding-carts", Stage: "dev"}

	nextStage, err := getNextStage(repo, data)
	assert.Equal(t, err, nil, "Received unexpected error")

	assert.Equal(t, nextStage, "staging", "Received unexpected stage")
}
