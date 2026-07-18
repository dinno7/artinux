package domain

import (
	"errors"
	"fmt"
)

var (

	ErrInvalidConfiguration = NewDomainError("invalid_configuration", "configuration is not valid")
	ErrNotFound             = NewDomainError("not_found", "resource not found")
	ErrInvalidInput         = NewDomainError("invalid_input", "invalid input provided")
	ErrConflict             = NewDomainError("conflict", "resource conflict")
	ErrInternal             = NewDomainError("internal_error", "an unexpected error occurred")
)

// DomainError represents a domain-specific error
// It keeps error details within the domain layer without leaking
// infrastructure concerns (HTTP status codes, DB specifics, etc.).
type DomainError struct {
	Code    string
	message string
	err     error
}

// NewDomainError creates a new DomainError.
func NewDomainError(code, message string) *DomainError {
	return &DomainError{
		Code:    code,
		message: message,
	}
}

// Wrap wraps an existing error with an existing DomainError.
func (e *DomainError) Wrap(err error) *DomainError {
	e.err = err
	return e
}

// Error implements the standard error interface.
func (e *DomainError) Error() string {
	if e.err != nil {
		return fmt.Errorf("%s: %s (caused by: %w)", e.Code, e.message, e.err).Error()
	}
	return fmt.Sprintf("%s: %s", e.Code, e.message)
}

// Message return pure error's message
func (de *DomainError) Message() string {
	return de.message
}

// MessageF formats a new message of existing Domain error.
func (de *DomainError) MessageF(format string, a ...any) *DomainError {
	de.message = fmt.Sprintf(format, a...)
	return de
}

// Unwrap allows errors.Is and errors.As to work with wrapped errors.
func (de *DomainError) UnWrap() error {
	return de.err
}

// Is allows checking error codes with errors.Is when comparing against sentinel errors.
func (e *DomainError) Is(target error) bool {
	t, ok := target.(*DomainError)
	if !ok {
		return errors.Is(e.err, target)
	}
	return e.Code == t.Code
}
