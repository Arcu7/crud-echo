package models

import (
	"errors"
)

const (
	InternalServerError  = "internal server error"
	BadRequest           = "bad request"
	NotFound             = "record not found"
	EmptyTable           = "table is empty"
	ResourceAlreadyExist = "resource already exist"
	InvalidParam         = "invalid parameter"
	ValidationError      = "validation error"
)

var (
	ErrInternalServerError  = errors.New("internal server error")
	ErrBadRequest           = errors.New("bad request")
	ErrNotFound             = errors.New("record not found")
	ErrRecordNotCreated     = errors.New("record not created")
	ErrEmptyTable           = errors.New("table is empty")
	ErrResourceAlreadyExist = errors.New("resource already exist")
	ErrInvalidParam         = errors.New("invalid parameter")
	ErrValidationError      = errors.New("validation error")
)

func GetErrorHTTPStatusCode(err error) int {
	switch {
	case errors.Is(err, ErrEmptyTable):
		return 200
	case errors.Is(err, ErrInternalServerError):
		return 500
	case errors.Is(err, ErrBadRequest), errors.Is(err, ErrInvalidParam):
		return 400
	case errors.Is(err, ErrNotFound):
		return 404
	case errors.Is(err, ErrResourceAlreadyExist):
		return 409
	case errors.Is(err, ErrValidationError):
		return 422
	default:
		return 500
	}
}

// for generic/human-readable message
func GetErrorHTTPStatusMessage(err error) string {
	switch {
	case errors.Is(err, ErrInternalServerError):
		return InternalServerError
	case errors.Is(err, ErrBadRequest):
		return BadRequest
	case errors.Is(err, ErrNotFound):
		return NotFound
	case errors.Is(err, ErrEmptyTable):
		return EmptyTable
	case errors.Is(err, ErrResourceAlreadyExist):
		return ResourceAlreadyExist
	case errors.Is(err, ErrInvalidParam):
		return InvalidParam
	case errors.Is(err, ErrValidationError):
		return ValidationError
	default:
		return InternalServerError
	}
}
