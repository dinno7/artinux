package usecases

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/dinno7/artinux/internal/domain"
	"github.com/dinno7/artinux/internal/domain/entities"
	"github.com/dinno7/artinux/internal/domain/ports"
)

type DownloadArtifactUC struct {
	storage ports.ObjectStorage
	hasher  ports.ChecksumHasher
	logger  ports.Logger
}

func NewDownloadArtifactUC(
	logger ports.Logger,
	storage ports.ObjectStorage,
	hasher ports.ChecksumHasher,
) *DownloadArtifactUC {
	return &DownloadArtifactUC{
		logger:  logger.With("usecase", "download_artifact"),
		storage: storage,
		hasher:  hasher,
	}
}

type DownloadArtifactInput struct {
	ObjectKey  string
	OutputPath string
}

func (uc *DownloadArtifactUC) Execute(
	ctx context.Context,
	input DownloadArtifactInput,
) (*entities.Artifact, error) {
	uc.logger.Info("Getting reader from storage", "object_key", input.ObjectKey)
	reader, artifact, err := uc.storage.Download(ctx, input.ObjectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to %w", err)
	}

	uc.logger.Info("Preparing file for writing to it", "path", input.OutputPath)
	file, err := os.OpenFile(input.OutputPath, os.O_CREATE|os.O_RDWR, 0o744)
	if err != nil {
		uc.logger.Error(
			"Failed to opening local file",
			err,
			"object_key", input.ObjectKey,
			"path", input.OutputPath,
		)
		return nil, domain.ErrInternal.Wrap(err)
	}

	uc.logger.Info("Writing to file & compute hash")
	teeReader := io.TeeReader(reader, uc.hasher.AsWriter())
	buf := make([]byte, 10*1024*1024) // 10mb
	_, err = io.CopyBuffer(file, teeReader, buf)
	if err != nil {
		uc.logger.Error(
			"Failed to writing in local file",
			err,
			"object_key", input.ObjectKey,
			"path", input.OutputPath,
		)
		return nil, domain.ErrInternal.Wrap(err)
	}

	uc.logger.Info("Checking checksum")
	hashedBase64 := uc.hasher.ComputeToBase64()
	if hashedBase64 != artifact.Checksum {
		if err := os.Remove(input.OutputPath); err != nil {
			uc.logger.Warn(
				"Failed to removing spoiled file",
				"error", err,
				"object_key", input.ObjectKey,
				"path", input.OutputPath,
			)
		}

		uc.logger.Error(
			"Downloaded file has not same checksum as storage, file corrupted",
			nil,
			"object_key", input.ObjectKey,
			"path", input.OutputPath,
		)
		return nil, domain.ErrFileChecksumNotSame
	}

	return artifact, nil
}
