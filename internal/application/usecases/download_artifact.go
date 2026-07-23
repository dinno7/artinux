package usecases

import (
	"context"
	"fmt"
	"io"

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
	ObjectKey string
}

type DownloadArtifactOutput struct {
	Artifact   *entities.Artifact
	FileReader io.ReadSeekCloser
}

func (uc *DownloadArtifactUC) Execute(
	ctx context.Context,
	input DownloadArtifactInput,
) (*DownloadArtifactOutput, error) {
	uc.logger.Info("Getting reader from storage", "object_key", input.ObjectKey)
	reader, artifact, err := uc.storage.Download(ctx, input.ObjectKey)
	if err != nil {
		return nil, err
	}

	uc.logger.Info("Checking checksum")
	hasValidChecksum, err := uc.hasher.CompareFromReader(reader, artifact.Checksum)
	if err != nil {
		return nil, fmt.Errorf("failed to %w", err)
	}

	if !hasValidChecksum {
		uc.logger.Error(
			"Downloaded file has not same checksum as storage, file corrupted",
			nil,
			"object_key", input.ObjectKey,
		)
		return nil, domain.ErrFileChecksumNotSame
	}

	// NOTE: for use again
	if _, err := reader.Seek(0, io.SeekStart); err != nil {
		uc.logger.Error(
			"Failed to reset opened file reader's pointer",
			err,
			"object_key", input.ObjectKey,
		)
		return nil, domain.ErrInternal.Wrap(err)
	}

	return &DownloadArtifactOutput{
		Artifact:   artifact,
		FileReader: reader,
	}, nil
}
