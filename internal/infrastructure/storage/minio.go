package storage

import (
	"cmp"
	"context"
	"errors"
	"io"
	"math"
	"strings"
	"time"

	"github.com/dinno7/artinux/internal/domain"
	"github.com/dinno7/artinux/internal/domain/entities"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minioStorage struct {
	client     *minio.Client
	bucketName string
	region     string
}

type MinIOConfig struct {
	Endpoint         string
	AccessKeyID      string
	SecretAccessKey  string
	BucketName       string
	Region           string
	HealthInterval   time.Duration
	UseSSL           bool
	MaxUploadRetries int
}

func NewMinIOStorage(cfg MinIOConfig) (*minioStorage, error) {
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:           credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure:          cfg.UseSSL,
		MaxRetries:      cmp.Or(cfg.MaxUploadRetries, 10),
		Region:          cfg.Region,
		TrailingHeaders: true,
	})
	if err != nil {
		return nil, err
	}

	s := &minioStorage{
		client:     minioClient,
		bucketName: cfg.BucketName,
		region:     cfg.Region,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err := s.Ping(ctx); err != nil {
		return nil, err
	}

	if err := s.CreateBucket(ctx, cfg.BucketName, cfg.Region); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *minioStorage) Name() string {
	return "MinIO S3 Object Storage"
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
		return domain.ErrStorageBucketExists.Wrap(err)
	}
	if !isBucketExists {
		if err := s.client.MakeBucket(
			ctx,
			normBucketName,
			minio.MakeBucketOptions{Region: region, ObjectLocking: true, ForceCreate: false},
		); err != nil {
			return domain.ErrStorageFailedCreateBucket.Wrap(err)
		}
	}
	return nil
}

func (s *minioStorage) ClearBucket(ctx context.Context) error {
	objs := s.client.ListObjects(ctx, s.bucketName, minio.ListObjectsOptions{
		WithMetadata: false,
		Recursive:    true,
	})

	objectKeys := []string{}
	for obj := range objs {
		objectKeys = append(objectKeys, obj.Key)
	}
	if len(objectKeys) > 0 {
		if err := s.DeleteBatch(ctx, objectKeys); err != nil {
			return err
		}
	}
	return nil
}

func (s *minioStorage) Upload(
	ctx context.Context,
	reader io.Reader,
	artifact *entities.Artifact,
) (string, error) {
	info, err := s.client.PutObject(
		ctx,
		s.bucketName,
		artifact.ObjectKey,
		reader,
		artifact.Size,
		minio.PutObjectOptions{
			UserMetadata: artifact.ToMap(),
			Checksum:     minio.ChecksumSHA256,
			ContentType:  "application/octet-stream",
		},
	)
	if err != nil {
		return "", domain.ErrStorageFailedToUpload.Wrap(err)
	}

	return info.ChecksumSHA256, nil
}

func (s *minioStorage) Download(
	ctx context.Context,
	objectKey string,
) (io.ReadSeekCloser, *entities.Artifact, error) {
	obj, err := s.client.GetObject(ctx, s.bucketName, objectKey, minio.GetObjectOptions{
		Checksum: true,
	})
	if err != nil {
		return nil, nil, domain.ErrStorageFailedToDownload.Wrap(err)
	}
	stats, err := obj.Stat()
	if err != nil {
		return nil, nil, domain.ErrStorageFailedToGetMetadata.Wrap(err)
	}

	artifact := ObjectStorageToArtifact(&stats)

	return obj, artifact, nil
}

func (s *minioStorage) ListObjects(
	ctx context.Context,
	prefix string,
	limit int,
) ([]*entities.Artifact, error) {
	list := s.client.ListObjects(ctx, s.bucketName, minio.ListObjectsOptions{
		WithMetadata: true,
		MaxKeys:      limit,
		WithVersions: true,
		Recursive:    true,
		Prefix:       prefix,
	})
	artifacts := make([]*entities.Artifact, int(math.Min(float64(len(list)), float64(limit))))

	for obj := range list {
		if len(artifacts) >= limit {
			break
		}
		a := ObjectStorageToArtifact(&obj)
		artifacts = append(artifacts, a)
	}

	return artifacts, nil
}

func (s *minioStorage) DeleteObject(ctx context.Context, objectKey string) error {
	if err := s.client.RemoveObject(
		ctx,
		s.bucketName,
		objectKey,
		minio.RemoveObjectOptions{},
	); err != nil {
		return domain.ErrStorageFailedToDeleteObject.Wrap(err)
	}
	return nil
}

func (s *minioStorage) DeleteBatch(ctx context.Context, objectKeys []string) error {
	objectsCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objectsCh)
		for _, k := range objectKeys {
			objectsCh <- minio.ObjectInfo{Key: k}
		}
	}()

	errs := s.client.RemoveObjects(ctx, s.bucketName, objectsCh, minio.RemoveObjectsOptions{})

	var err error = nil
	for delErr := range errs {
		err = errors.Join(delErr.Err)
	}

	return err
}

