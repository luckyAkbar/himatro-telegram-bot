package usecase

import "errors"

var (
	ErrValidation         = errors.New("validation error")
	ErrInternal           = errors.New("internal error")
	ErrExternalService    = errors.New("external service error")
	ErrExternalBadRequest = errors.New("external return bad request")
	ErrNotFound           = errors.New("not found ")
	ErrUnauthorized       = errors.New("unauthorized request")
)
