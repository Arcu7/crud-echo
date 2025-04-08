package models

import (
	"errors"
)

const (
	InternalServerError = "internal server error"
	BadRequest          = "bad request"
	NotFound            = "record not found"
)

var (
	ErrInternalServerError  = errors.New("internal server error")
	ErrBadRequest           = errors.New("bad request")
	ErrNotFound             = errors.New("record not found")
	ErrTableEmpty           = errors.New("table is empty")
	ErrResourceExistAlready = errors.New("resource already exists")
	ErrInvalidParam         = errors.New("invalid parameter")
)

type ValidationError struct {
	Message string
	Errors  map[string]string
}

func (e *ValidationError) Error() string {
	return e.Message
}
