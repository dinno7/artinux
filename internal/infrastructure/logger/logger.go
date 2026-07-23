package logger

import (
	"sync"

	"github.com/dinno7/artinux/internal/domain/ports"
)

var (
	once   sync.Once
	global ports.Logger
)

type LoggerConfig struct {
	Format string
	Level  string
}

func NewLogger(loggerConfig LoggerConfig) ports.Logger {
	once.Do(func() {
		global = newZapLogger(&loggerConfig)
	})
	return global
}
