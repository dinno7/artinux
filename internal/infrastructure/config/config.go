package config

import "time"

type Config struct {
	Env           string        `mapstructure:"env"`
	ObjectStorage ObjectStorage `mapstructure:"object_storage" validate:"required"`
	Upload        Upload        `mapstructure:"upload"`
	Logging       Logging       `mapstructure:"logging"        validate:"required"`
}

type ObjectStorage struct {
	Endpoint            string        `mapstructure:"endpoint"`
	Username            string        `mapstructure:"username"`
	Password            string        `mapstructure:"password"`
	Region              string        `mapstructure:"region"`
	UseSSL              bool          `mapstructure:"use_ssl"`
	BucketName          string        `mapstructure:"bucket_name"`
	HealthCheckInterval time.Duration `mapstructure:"health_check_interval"`
}

type Upload struct {
	MaxSizeMB             int64    `mapstructure:"max_size_mb"       validate:"required"`
	AllowedFileExtensions []string `mapstructure:"allowed_file_exts" validate:"required"`
}

type Logging struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}
