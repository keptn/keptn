package common

import (
	"errors"
	"fmt"
)

var ErrConfigStoreInvalidToken = errors.New("invalid git token")

var ErrConfigStoreUpstreamNotFound = errors.New("upstream repository not found")

var ErrSequenceWithTriggeredIDAlreadyExists = errors.New("sequence with the same triggeredID already exists")

var ErrOpenRemediationNotFound = errors.New("open remediation not found")

var ErrProjectAlreadyExists = errors.New("project already exists")

var ErrServiceAlreadyExists = errors.New("service already exists")

var ErrServiceNotFound = errors.New("service not found")

var ErrProjectNotFound = errors.New("project not found")

var ErrInvalidStageChange = errors.New("stage name cannot be changed or removed")

var ErrStageNotFound = errors.New("stage not found")

var ErrChangesRollback = errors.New("failed to rollback changes")

var ErrSequencePaused = errors.New("sequence is paused")

var ErrSequenceBlocked = errors.New("sequence is currently blocked")

var ErrSequenceBlockedWaiting = errors.New("sequence is currently blocked by waiting for another sequence to end")

var ErrNoMatchingEvent = errors.New("no matching event found")

var ErrSequenceNotFound = errors.New("sequence not found")

var ErrInternalError = errors.New("internal server error")

var InvalidRequestFormatMsg = "Invalid request format: %s"

var UnexpectedErrorFormatMsg = "Unexpected error: %s"

var UnableRetrieveLogsMsg = "Unable to retrieve logs: %s"

var ProjectNotFoundMsg = "Project not found: %s"

var SequenceNotFoundMsg = "Sequence not found: %s"

var EventNotFoundMsg = "Event not found: %s"

var InvalidPayloadMsg = "Could not validate payload: %s"

var NoProjectNameMsg = "Must provide a project name"

var NoServiceNameMsg = "Must provide a service name"

var UnableQueryStateMsg = "Unable to query sequence state repository: %s"

var UnableQuerySequenceExecutionMsg = "Unable to query sequence execution repository: %s"

var UnableControleSequenceMsg = "Unable to control sequence: %s"

var UnableFindSequenceMsg = "Unable to control sequence: %s"

var InvalidRemoteURLMsg = "Invalid RemoteURL: %s"

var UnableQueryIntegrationsMsg = "Unable to query uniform integrations repository: %s"

var UnableMarshallProvisioningData = "Error marshalling provisioning data: %s"

var UnableUnMarshallProvisioningData = "Error unmarshalling provisioning data: %s"

var UnableReadProvisioningData = "Error reading provisioning data: %s"

var UnableProvisionInstance = "Error provisioning a project instance: %s"

var UnableProvisionInstanceGeneric = fmt.Sprintf(UnableProvisionInstance, "unable to provision an instance")

var UnableProvisionDelete = "Error deleting a provisioned project instance: %s"

var UnableProvisionDeleteGeneric = fmt.Sprintf(UnableProvisionDelete, "unable to delete an instance")

var UnableProvisionDeleteReq = "Error creating delete provision request: %s"

var UnableProvisionPostReq = "Error creating post provision request: %s"

var OtherActiveSequencesRunning = "Other sequences are currently running in the same stage for the same service with context id: "
