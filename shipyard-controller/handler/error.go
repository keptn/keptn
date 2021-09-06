package handler

import "errors"

var ErrProjectAlreadyExists = errors.New("project already exists")

var ErrServiceAlreadyExists = errors.New("service already exists")

var ErrServiceNotFound = errors.New("service not found")

var ErrProjectNotFound = errors.New("project not found")

var ErrStageNotFound = errors.New("stage not found")
