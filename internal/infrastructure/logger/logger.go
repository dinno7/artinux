package logger

import (
	"sync"

	"github.com/dinno7/artinux/internal/domain/ports"
)

var once sync.Once

type LoggerConfig struct {
	Format string
	Level  string
}

func NewLogger(loggerConfig LoggerConfig) ports.Logger {
	var l ports.Logger
	once.Do(func() {
		l = newZapLogger(&loggerConfig)
	})
	return l
}
