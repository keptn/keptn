package handler

import "errors"

var ErrProjectAlreadyExists = errors.New("project already exists")

var errServiceAlreadyExists = errors.New("project already exists")

var errServiceNotFound = errors.New("service not found")

var errProjectNotFound = errors.New("project not found")

var errStageNotFound = errors.New("stage not found")
