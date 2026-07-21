package config

import (
	"errors"
	"regexp"
	"slices"
)

var (
	ErrInvalidEnv = errors.New("invalid 'env' value, can be of of dev, prod")

	ErrInvalidLogLevel = errors.New(
		"invalid 'logging.level' value, can be of of debug, info, warn, error, fatal",
	)
	ErrInvalidLogFormat = errors.New(
		"invalid 'logging.format' value, can be of of json, text",
	)

	ErrRequiredStorageEndpoint   = errors.New("'object_storage.endpoint' is required")
	ErrRequiredStorageBucketName = errors.New("'object_storage.bucket_name' is required")
	ErrRequiredStorageUsername   = errors.New("'object_storage.username' is required")
	ErrRequiredStoragePassword   = errors.New("'object_storage.password' is required")
	ErrRequiredStorageRegion     = errors.New("'object_storage.region' is required")
	ErrInvalidStorageBucketName  = errors.New(
		"'object_storage.bucket_name' can not contains upper letter",
	)

	ErrRequiredUploadMaxSize         = errors.New("'upload.max_size_mb' is required")
	ErrRequiredUploadAllowedFileExts = errors.New(
		"'upload.allowed_file_exts' should have at least 1 ext",
	)
)

var (
	validEnv       = []string{"dev", "prod"}
	validLogLevel  = []string{"debug", "info", "warn", "error", "fatal"}
	validLogFormat = []string{"json", "text"}
)

func (cfg *Config) validate() error {
	// INFO: Env cfg
	if !slices.Contains(validEnv, cfg.Env) {
		return ErrInvalidEnv
	}

	// INFO: Logging cfg
	if !slices.Contains(validLogLevel, cfg.Logging.Level) {
		return ErrInvalidLogLevel
	}
	if !slices.Contains(validLogFormat, cfg.Logging.Format) {
		return ErrInvalidLogLevel
	}

	// INFO: Object storage cfg
	if cfg.ObjectStorage.Endpoint == "" {
		return ErrRequiredStorageEndpoint
	}
	if cfg.ObjectStorage.BucketName == "" {
		return ErrRequiredStorageBucketName
	}
	bucketRX := regexp.MustCompile(`[A-Z]`)
	if bucketRX.MatchString(cfg.ObjectStorage.BucketName) {
		return ErrInvalidStorageBucketName
	}
	if cfg.ObjectStorage.Username == "" {
		return ErrRequiredStorageUsername
	}
	if cfg.ObjectStorage.Password == "" {
		return ErrRequiredStoragePassword
	}
	if cfg.ObjectStorage.Region == "" {
		return ErrRequiredStorageRegion
	}

	// INFO: Upload cfg
	if cfg.Upload.MaxSizeMB < 0 {
		return ErrRequiredUploadMaxSize
	}
	if len(cfg.Upload.AllowedFileExtensions) == 0 {
		return ErrRequiredUploadAllowedFileExts
	}

	return nil
}
