package config

import (
	"testing"
)

func TestConfigValidate(t *testing.T) {
	validConfig := Config{
		Env: "dev",
		Logging: Logging{
			Level:  "debug",
			Format: "text",
		},
		ObjectStorage: ObjectStorage{
			Endpoint:   "localhost:9000",
			Username:   "admin",
			Password:   "secret",
			BucketName: "artifacts",
			Region:     "us-east-1",
		},
		Upload: Upload{
			MaxSizeMB:             100,
			AllowedFileExtensions: []string{"deb", "rpm"},
		},
		HTTPServer: HTTPServer{
			Schema: "http",
			Host:   "0.0.0.0",
			Port:   7000,
		},
	}

	t.Run("valid config passes", func(t *testing.T) {
		err := validConfig.validate()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	tests := []struct {
		name  string
		apply func(*Config)
		want  error
	}{
		{
			name:  "invalid env",
			apply: func(c *Config) { c.Env = "staging" },
			want:  ErrInvalidEnv,
		},
		{
			name:  "invalid log level",
			apply: func(c *Config) { c.Logging.Level = "trace" },
			want:  ErrInvalidLogLevel,
		},
		{
			name:  "invalid log format",
			apply: func(c *Config) { c.Logging.Format = "xml" },
			want:  ErrInvalidLogFormat,
		},
		{
			name:  "empty endpoint",
			apply: func(c *Config) { c.ObjectStorage.Endpoint = "" },
			want:  ErrRequiredStorageEndpoint,
		},
		{
			name:  "empty bucket name",
			apply: func(c *Config) { c.ObjectStorage.BucketName = "" },
			want:  ErrRequiredStorageBucketName,
		},
		{
			name:  "bucket name with uppercase",
			apply: func(c *Config) { c.ObjectStorage.BucketName = "MyBucket" },
			want:  ErrInvalidStorageBucketName,
		},
		{
			name:  "empty username",
			apply: func(c *Config) { c.ObjectStorage.Username = "" },
			want:  ErrRequiredStorageUsername,
		},
		{
			name:  "empty password",
			apply: func(c *Config) { c.ObjectStorage.Password = "" },
			want:  ErrRequiredStoragePassword,
		},
		{
			name:  "empty region",
			apply: func(c *Config) { c.ObjectStorage.Region = "" },
			want:  ErrRequiredStorageRegion,
		},
		{
			name:  "negative max size",
			apply: func(c *Config) { c.Upload.MaxSizeMB = -1 },
			want:  ErrRequiredUploadMaxSize,
		},
		{
			name:  "empty allowed extensions",
			apply: func(c *Config) { c.Upload.AllowedFileExtensions = []string{} },
			want:  ErrRequiredUploadAllowedFileExts,
		},
		{
			name:  "extension starts with dot",
			apply: func(c *Config) { c.Upload.AllowedFileExtensions = []string{".deb"} },
			want:  ErrFileExtStartsWithDot,
		},
		{
			name:  "port too low",
			apply: func(c *Config) { c.HTTPServer.Port = -1 },
			want:  ErrHTTPServerInvalidPort,
		},
		{
			name:  "port too high",
			apply: func(c *Config) { c.HTTPServer.Port = 70000 },
			want:  ErrHTTPServerInvalidPort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := validConfig
			tt.apply(&cfg)
			err := cfg.validate()
			if err != tt.want {
				t.Errorf("expected %v, got %v", tt.want, err)
			}
		})
	}
}
