package errors

import (
	"reflect"
)

type ResourceServiceError struct {
	s string
}

func New(s string) *ResourceServiceError {
	return &ResourceServiceError{s}
}
func (e *ResourceServiceError) Is(target error) bool {
	return reflect.DeepEqual(e.s, target.Error()) || target.Error() == ""
}

func (e *ResourceServiceError) Error() string {
	return e.s
}

// Project specific errors

var ErrProjectNotFound = New("project not found")
var ErrProjectAlreadyExists = New("project already exists")

// Stage specific errors

var ErrStageNotFound = New("stage not found")
var ErrStageAlreadyExists = New("stage already exists")

// Service Specific errors

var ErrServiceNotFound = New("service not found")
var ErrServiceAlreadyExists = New("service already exists")

// Resource specific errors

var ErrResourceNotFound = New("resource not found")
var ErrResourceAlreadyExists = New("resource already exists")
var ErrResourceNotBase64Encoded = New("resource content is not base64 encoded")
var ErrResourceInvalidResourceURI = New("invalid resource uri")

// Git specific errors

var ErrInvalidGitToken = New("invalid git token")
var ErrRepositoryNotFound = New("repository not found")
var ErrInvalidGitContext = New("invalid git context")
var ErrResolvedNilHash = New("resolved nil hash")
var ErrResolveRevision = New("revision does not exist")
var ErrBranchExists = New("branch already exists")
var ErrBranchNotFound = New("branch not found")
var ErrTagExists = New("tag already exists")
var ErrTagNotFound = New("tag not found")
var ErrAnonymousRemoteName = New("anonymous remote name must be 'anonymous'")
var ErrWorktreeNotProvided = New("worktree should be provided")
var ErrEmptyRemoteRepository = New("remote repository is empty")
var ErrAuthenticationRequired = New("authentication required")
var ErrAuthorizationFailed = New("authorization failed")
var ErrEmptyUploadPackRequest = New("empty git-upload-pack given")
var ErrInvalidAuthMethod = New("invalid auth method")
var ErrAlreadyConnected = New("session already established")

//Git worktree
var ErrWorktreeNotClean = New("worktree is not clean")
var ErrSubmoduleNotFound = New("submodule not found")
var ErrUnstagedChanges = New("worktree contains unstaged changes")
var ErrGitModulesSymlink = New(".gitmodules is a symlink")
var ErrNonFastForwardUpdate = New("non-fast-forward update")

// Git repo
var ErrInvalidReference = New("invalid reference, should be a tag or a branch")
var ErrRepositoryNotExists = New("repository does not exist")
var ErrRepositoryIncomplete = New("repository's commondir path does not exist")
var ErrRepositoryAlreadyExists = New("repository already exists")
var ErrRemoteNotFound = New("remote not found")
var ErrRemoteExists = New("remote already exists")
var ErrIsBareRepository = New("worktree not available in a bare repository")
var ErrUnableToResolveCommit = New("unable to resolve commit")
var ErrPackedObjectsNotSupported = New("packed objects not supported")
var ErrReferenceNotFound = New("reference not found")
var NoErrAlreadyUpToDate = New("already up-to-date")
var ErrDeleteRefNotSupported = New("server does not support delete-refs")
var ErrForceNeeded = New("some refs were not updated")
var ErrExactSHA1NotSupported = New("server does not support exact SHA1 refspec")

// Credential specific errors

var ErrCredentialsNotFound = New("could not find upstream repository credentials")
var ErrMalformedCredentials = New("could not decode upstream repository credentials")
var ErrCredentialsInvalidRemoteURL = New("invalid remote URL")
var ErrCredentialsTokenMustNotBeEmpty = New("token must not be empty")
var ErrCredentialsPrivateKeyMustNotBeEmpty = New("private key must not be empty")
var ErrProxyInvalidScheme = New("proxy scheme must be http or https")
var ErrInvalidRemoteURL = New("RemoteURL scheme must be http, https or ssh")
var ErrProxyInvalidURL = New("proxy URL must contain IP address and port (<ip-address>:<port>)")

// Error messages

const ErrMsgCouldNotRetrieveCredentials = "could not read credentials for project %s: %w"
const ErrMsgInvalidRequestFormat = "Invalid request format"
const ErrMsgCouldNotSetUser = "could not set git user: %w"
const ErrMsgCouldNotCreatePath = "could not create path %s: %w"
const ErrMsgCouldNotGitAction = "could not %s git repo for project %s: %w"
const ErrMsgCouldNotCommit = "could not commit changes in project %s: %w"
const ErrMsgCouldNotGetRevision = "could not get current revision for project %s: %w"
const ErrMsgCouldNotGetDefBranch = "could not get default branch for project %s: %w"
const ErrMsgCouldNotCheckout = "could not checkout branch %s: %w"
const ErrMsgCouldNotCreate = "could not create branch %s for project %s: %w"
