package cmd

import (
	"errors"
	"fmt"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
)

type sequenceControlStruct struct {
	keptnContext *string
	project      *string
	stage        *string
}

type SequenceState string

const (
	pauseSequence  SequenceState = "pause"
	resumeSequence SequenceState = "resume"
	abortSequence  SequenceState = "abort"
)

func AbortSequence(params sequenceControlStruct) error {
	return controlSequence(abortSequence, params)
}

func PauseSequence(params sequenceControlStruct) error {
	return controlSequence(pauseSequence, params)
}

func ResumeSequence(params sequenceControlStruct) error {
	return controlSequence(resumeSequence, params)
}

func controlSequence(sequenceState SequenceState, params sequenceControlStruct) error {
	endPoint, apiToken, err := credentialmanager.NewCredentialManager(assumeYes).GetCreds(namespace)
	if err != nil {
		return errors.New(authErrorMsg)
	}
	if endPointErr := CheckEndpointStatus(endPoint.String()); endPointErr != nil {
		return fmt.Errorf("Error connecting to server: %s"+endPointErrorReasons,
			endPointErr)
	}
	sequenceControlHandler := apiutils.NewAuthenticatedSequenceControlHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
	controlParams := apiutils.SequenceControlParams{
		Project:      *params.project,
		KeptnContext: *params.keptnContext,
		Stage:        *params.stage,
		State:        string(sequenceState),
	}

	return sequenceControlHandler.ControlSequence(controlParams)
}
