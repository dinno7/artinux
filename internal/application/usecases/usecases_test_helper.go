package usecases

import (
	"github.com/dinno7/artinux/internal/domain/ports"
)

// noopLogger is a silent logger for use in tests.
type noopLogger struct{}

func (n *noopLogger) Debug(_ string, _ ...any) {}
func (n *noopLogger) Info(_ string, _ ...any)  {}
func (n *noopLogger) Warn(_ string, _ ...any)  {}
func (n *noopLogger) Error(_ string, _ error, _ ...any) {
}
func (n *noopLogger) Fatal(_ string, _ error, _ ...any) {
}
func (n *noopLogger) With(_ ...any) ports.Logger { return n }
