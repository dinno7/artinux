package storage

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/dinno7/artinux/internal/domain"
	"github.com/dinno7/artinux/internal/domain/ports"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minioStorage struct {
	client *minio.Client
}

type MinIOConfig struct {
	Endpoint            string
	AccessKeyID         string
	SecretAccessKey     string
	HealthCheckInterval time.Duration
	UseSSL              bool
}

func NewMinIOStorage(cfg MinIOConfig) (ports.ObjectStorage, error) {
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	s := &minioStorage{
		client: minioClient,
	}

	if err := s.Ping(context.Background()); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *minioStorage) Ping(ctx context.Context) error {
	if _, err := s.client.ListBuckets(ctx); err != nil {
		return domain.ErrStorageUnavailable.Wrap(err)
	}
	return nil
}

func (s *minioStorage) CreateBucket(ctx context.Context, bucketName string, region string) error {
	normBucketName := strings.TrimSpace(strings.ToLower(bucketName))

	isBucketExists, err := s.client.BucketExists(ctx, normBucketName)
	if err != nil {
		return err
	}
	if !isBucketExists {
		return s.client.MakeBucket(
			ctx,
			normBucketName,
			minio.MakeBucketOptions{Region: region, ObjectLocking: true, ForceCreate: false},
		)
	}
	return nil
}
