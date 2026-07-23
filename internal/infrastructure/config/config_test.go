package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func getTestConfig() Config {
	return Config{
		Env: "dev",
		Upload: Upload{
			MaxSizeMB:             10,
			AllowedFileExtensions: []string{},
		},
		ObjectStorage: ObjectStorage{
			Endpoint:   "localhost:9090",
			Username:   "admin",
			Password:   "someunsafepassword",
			BucketName: "bucket_name",
			Region:     "us-east-1",
		},
		Logging: Logging{
			Level:  "debug",
			Format: "text",
		},
	}
}

func TestIsProduction(t *testing.T) {
	cfg := getTestConfig()

	cfg.Env = "prod"
	isProd := cfg.IsProduction()

	if !isProd {
		t.Errorf("Expected IsProduction return true when env is prod, but got %v", isProd)
	}

	cfg.Env = "dev"
	isProd = cfg.IsProduction()
	if isProd {
		t.Errorf("Expected IsProduction return false when env is not prod, but got %v", isProd)
	}
}

func TestValidationUpperBucketName(t *testing.T) {
	cfg := getTestConfig()
	cfg.ObjectStorage.BucketName = "BucketName"
	err := cfg.validate()
	require.ErrorIs(t, err, ErrInvalidStorageBucketName)
}
