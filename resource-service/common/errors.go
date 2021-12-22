package common

import "errors"

// Project specific errors

var ErrProjectNotFound = errors.New("project not found")
var ErrProjectAlreadyExists = errors.New("project already exists")

// Stage specific errors

var ErrStageNotFound = errors.New("stage not found")
var ErrStageAlreadyExists = errors.New("stage already exists")

// Service Specific errors

var ErrServiceNotFound = errors.New("project not found")
var ErrServiceAlreadyExists = errors.New("project already exists")

// Resource specific errors

var ErrResourceNotFound = errors.New("resource not found")
var ErrResourceAlreadyExists = errors.New("resource already exists")

// Git specific errors

var ErrInvalidGitToken = errors.New("invalid git token")
var ErrRepositoryNotFound = errors.New("upstream repository not found")
