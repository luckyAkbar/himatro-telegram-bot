package himatro

import "errors"

var (
	ErrValidation      = errors.New("validation error")
	ErrInternal        = errors.New("internal error")
	ErrExternalService = errors.New("external service error")
	ErrBadRequest      = errors.New("bad request")
	ErrForbidden       = errors.New("action is forbidden")
	ErrNotFound        = errors.New("not found")
	ErrUnauthorized    = errors.New("unauthorized")
)
