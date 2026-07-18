package ports

import (
	"context"
	"io"
)

type ObjectStorage interface {
	Ping(ctx context.Context) error
}
