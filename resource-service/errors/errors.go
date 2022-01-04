
package errors

import "errors"

// Project specific errors

var ErrProjectNotFound = errors.New("project not found")
var ErrProjectAlreadyExists = errors.New("project already exists")

// Stage specific errors

var ErrStageNotFound = errors.New("stage not found")
var ErrStageAlreadyExists = errors.New("stage already exists")

// Service Specific errors

var ErrServiceNotFound = errors.New("service not found")
var ErrServiceAlreadyExists = errors.New("service already exists")

// Resource specific errors

var ErrResourceNotFound = errors.New("resource not found")
var ErrResourceAlreadyExists = errors.New("resource already exists")
var ErrResourceNotBase64Encoded = errors.New("resource content is not base64 encoded")
var ErrResourceInvalidResourceURI = errors.New("invalid resource uri")

// Git specific errors

var ErrInvalidGitToken = errors.New("invalid git token")
var ErrRepositoryNotFound = errors.New("upstream repository not found")
var ErrInvalidGitContext = errors.New("invalid git context")
var ErrResolvedNilHash = errors.New("resolved nil hash")
var ErrResolveRevision = errors.New("revision does not exist")

// Credential specific errors

var ErrCredentialsNotFound = errors.New("could not find upstream repository credentials")
var ErrMalformedCredentials = errors.New("could not decode upstream repository credentials")

// Error messages

const ErrMsgCouldNotRetrieveCredentials = "could not read credentials for project %s: %w"
const ErrMsgInvalidRequestFormat = "Invalid request format"
const ErrMsgCouldNotSetUser = "could not set git user: %w"
const ErrMsgCouldNotCreatePath = "could not create path %s: %w"
const ErrMsgCouldNotGitAction =  "could not %s git repo for project %s: %w"
const ErrMsgCouldNotCommit = "could not commit changes in project %s: %w"
const ErrMsgCouldNotGetRevision = "could not get current revision for project %s: %w"
const ErrMsgCouldNotGetDefBranch = "could not get default branch for project %s: %w"
const ErrMsgCouldNotCheckout = "could not checkout branch %s: %w"
const ErrMsgCouldNotCreate= "could not create branch %s for project %s: %w"
