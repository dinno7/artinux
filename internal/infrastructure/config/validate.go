package config

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"strings"
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
	ErrFileExtStartsWithDot = errors.New(
		"'upload.allowed_file_exts' should not starts with '.'",
	)
	ErrHTTPServerInvalidPort = errors.New(
		"http_server.port is not valid, the port must be between 0-65535",
	)
	ErrHTTPServerInvalidAddress = errors.New(
		"http_server.host, http_server.port invalid address for http server",
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
		return ErrInvalidLogFormat
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
	for _, ext := range cfg.Upload.AllowedFileExtensions {
		if strings.HasPrefix(ext, ".") {
			return ErrFileExtStartsWithDot
		}
	}

	// INFO: HTTP server
	if cfg.HTTPServer.Port < 0 || cfg.HTTPServer.Port > 65535 {
		return ErrHTTPServerInvalidPort
	}
	httpServerAddr := fmt.Sprintf(
		"%s://%s:%d",
		cfg.HTTPServer.Schema,
		cfg.HTTPServer.Host,
		cfg.HTTPServer.Port,
	)
	_, err := url.Parse(httpServerAddr)
	if err != nil {
		return errors.Join(ErrHTTPServerInvalidAddress, err)
	}

	return nil
}
