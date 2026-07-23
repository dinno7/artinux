package logger

import (
	"os"
	"strings"

	"github.com/dinno7/artinux/internal/domain/ports"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.SugaredLogger
}

func newZapLogger(cfg *LoggerConfig) *zapLogger {
	logLevel := resolveLogLevel(cfg.Level)
	writer := zapcore.AddSync(os.Stdout)
	var config zapcore.EncoderConfig
	var encoder zapcore.Encoder

	logFormat := strings.TrimSpace(strings.ToLower(cfg.Format))
	if logFormat == "json" {
		config = zap.NewProductionEncoderConfig()
		encoder = zapcore.NewJSONEncoder(config)
	} else {
		config = zap.NewDevelopmentEncoderConfig()
		config.ConsoleSeparator = " | "
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(config)
	}

	core := zapcore.NewCore(encoder, writer, logLevel)

	logger := &zapLogger{logger: zap.New(core).Sugar()}

	return logger
}

func (l *zapLogger) Debug(msg string, args ...any) {
	l.logger.Debugw(msg, args...)
}

func (l *zapLogger) Info(msg string, args ...any) {
	l.logger.Infow(msg, args...)
}

func (l *zapLogger) Warn(msg string, args ...any) {
	l.logger.Warnw(msg, args...)
}

func (l *zapLogger) Error(msg string, err error, args ...any) {
	if err != nil {
		args = append(args, "error", err)
	}
	l.logger.Errorw(msg, args...)
}

func (l *zapLogger) Fatal(msg string, err error, args ...any) {
	if err != nil {
		args = append(args, "error", err)
	}
	l.logger.Fatalw(msg, args...)
}

func (l *zapLogger) With(args ...any) ports.Logger {
	if len(args) == 0 {
		return l
	}
	// Validate args are key-value pairs
	if len(args)%2 != 0 {
		// Append a placeholder for odd args to prevent panic
		args = append(args, "(MISSING)")
	}
	return &zapLogger{
		logger: l.logger.With(args),
	}
}

func resolveLogLevel(strLevel string) zapcore.Level {
	logLevel := zapcore.InfoLevel
	switch strLevel {
	case "debug":
		logLevel = zapcore.DebugLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "fatal":
		logLevel = zapcore.FatalLevel
	}
	return logLevel
}
