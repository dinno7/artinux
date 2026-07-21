package config

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dinno7/artinux/internal/domain"
	"github.com/spf13/viper"
)

var getConfiguration = sync.OnceValues(func() (*Config, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// INFO: Name of the config file without an extension (Viper will intuit the type
	// from an extension on the actual file)
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(userHomeDir)

	// INFO: ENV override support
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

	// INFO: Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, domain.ErrInvalidConfiguration.Wrap(err)
	}

	configurations := getDefaultConfig()
	if err := viper.Unmarshal(configurations); err != nil {
		return nil, domain.ErrInvalidConfiguration.Wrap(err)
	}

	if err := configurations.validate(); err != nil {
		return nil, domain.ErrInvalidConfiguration.MessageF("configuration is not valid").Wrap(err)
	}

	return configurations, nil
})

func Get() (*Config, error) {
	return getConfiguration()
}

func getDefaultConfig() *Config {
	return &Config{
		Env: "prod",
		Upload: Upload{
			MaxSizeMB: 1024, // 1GB
			AllowedFileExtensions: []string{
				"AppImage",
				"deb",
				"rpm",
				"flatpak",
				"tar.gz", "tar.bz2", "tar.xz", "tgz", "gz",
				"zip",
				"xz",
				"bz2",
				"zst",
			},
		},
		Logging: Logging{
			Level:  "info",
			Format: "json",
		},
		ObjectStorage: ObjectStorage{
			BucketName:       "artifacts",
			Region:           "us-east-1",
			UseSSL:           false,
			MaxUploadRetries: 10,
			HealthInterval:   time.Second * 5,
		},
	}
}

func (cfg *Config) IsProduction() bool {
	return cfg.Env == "prod"
}
