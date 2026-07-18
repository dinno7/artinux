package ports

import (
	"context"
	"io"
)

type ObjectStorage interface {
	Ping(ctx context.Context) error
	CreateBucket(ctx context.Context, bucketName string, region string) error
}
