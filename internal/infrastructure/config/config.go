package config

import "time"

type Config struct {
	Env           string        `mapstructure:"env"`
	ObjectStorage ObjectStorage `mapstructure:"object_storage"`
	Upload        Upload        `mapstructure:"upload"`
	Logging       Logging       `mapstructure:"logging"       `
	HTTPServer    HTTPServer    `mapstructure:"http_server"`
}

type ObjectStorage struct {
	Endpoint         string        `mapstructure:"endpoint"`
	Username         string        `mapstructure:"username"`
	Password         string        `mapstructure:"password"`
	Region           string        `mapstructure:"region"`
	UseSSL           bool          `mapstructure:"use_ssl"`
	BucketName       string        `mapstructure:"bucket_name"`
	HealthInterval   time.Duration `mapstructure:"health_interval"`
	MaxUploadRetries int           `mapstructure:"max_retries"`
}

type Upload struct {
	MaxSizeMB             int64    `mapstructure:"max_size_mb"      `
	AllowedFileExtensions []string `mapstructure:"allowed_file_exts"`
}

type Logging struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type HTTPServer struct {
	Schema string `mapstructure:"schema"`
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
}
