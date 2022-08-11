package handler

import "errors"

var (
	ErrInternal        = errors.New("internal error")
	ErrExternalService = errors.New("external service return an error. may be caused invalid request")
	ErrBadRequest      = errors.New("bad request: please use correct command syntax. Get help by typing /command-name help")
	ErrNotFound        = errors.New("entities not found, please try again")
	ErrValidation      = errors.New("validation error. Try sending correct and valid data then try again")
	ErrUnauthorized    = errors.New("request unauthorized. Maybe try using login command first")
)
