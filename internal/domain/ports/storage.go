package ports

import (
	"context"
	"io"

	"github.com/dinno7/artinux/internal/domain/entities"
)

type ObjectStorage interface {
	Ping(ctx context.Context) error
	CreateBucket(ctx context.Context, bucketName string, region string) error
	ClearBucket(ctx context.Context) error
	Upload(ctx context.Context, reader io.Reader, artifact *entities.Artifact) (string, error)
	Download(ctx context.Context, objectKey string) (io.ReadSeekCloser, *entities.Artifact, error)
	ListObjects(ctx context.Context, prefix string, limit int) ([]*entities.Artifact, error)
	DeleteObject(ctx context.Context, objectKey string) error
	DeleteBatch(ctx context.Context, objectKeys []string) error
}
