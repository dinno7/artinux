package domain

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidExtension = NewDomainError(
		"invalid_file_extension",
		"file extension is not valid",
	)
	ErrFileEmptyChecksum = NewDomainError(
		"file_empty_checksum",
		"file checksum is empty",
	)
	ErrFileUploadNotSupportIntegrity = NewDomainError(
		"file_not_integrity",
		"file integrity violation",
	)
	ErrFileChecksumNotSame = NewDomainError(
		"file_invalid_checksum",
		"file's checksum is not the same as local one",
	)
	ErrFileTooLarge = NewDomainError("file_too_large", "file is too large")
	ErrPathNotFile  = NewDomainError(
		"not_a_file",
		"provided path is not a file",
	)
	ErrFileNotAccessible          = NewDomainError("file_not_accessible", "file not accessible")
	ErrFileFailedToDetectFileType = NewDomainError(
		"failed_detect_file_type",
		"file type's detection has been failed",
	)
	ErrFileIsEmpty          = NewDomainError("empty_file", "file is empty")
	ErrFileNotExists        = NewDomainError("file_not_exists", "file not exists")
	ErrAuthenticationFailed = NewDomainError("authentication_failed", "authentication failed")

	ErrStorageUnavailable = NewDomainError("storage_unavailable", "storage unavailable")
	ErrInvalidOsOrArch    = NewDomainError(
		"invalid_os_or_arch",
		"os's name or arch's name is not valid",
	)
	ErrStorageFailedToUpload = NewDomainError(
		"storage_failed_upload",
		"failed to upload object to storage",
	)
	ErrStorageFailedToDownload = NewDomainError(
		"storage_failed_download",
		"failed to download from object storage",
	)
	ErrStorageFailedToGetMetadata = NewDomainError(
		"storage_failed_metadata",
		"failed to getting object's metadata",
	)
	ErrStorageFailedToDeleteObject = NewDomainError(
		"storage_failed_delete_object",
		"failed to delete object with object key",
	)
	ErrStorageEmptyArtifactToUpload = NewDomainError(
		"empty_artifact_upload",
		"no artifact to upload",
	)
	ErrStorageFailedCreateBucket = NewDomainError(
		"failed_create_bucket",
		"failed to create new bucket",
	)
	ErrStorageFailedRemoveBucket = NewDomainError(
		"failed_remove_bucket",
		"failed to remove bucket",
	)
	ErrStorageBucketExists = NewDomainError(
		"bucket_exists",
		"bucket already exists",
	)

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
	if errors.Is(err, e) {
		return e
	}
	e.err = err
	return e
}

// Error implements the standard error interface.
func (e *DomainError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s (caused by: %s)", e.Code, e.message, e.err.Error())
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
