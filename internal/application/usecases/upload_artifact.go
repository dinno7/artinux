package usecases

import (
	"context"
	"io"
	"os"

	"github.com/dinno7/artinux/internal/domain"
	"github.com/dinno7/artinux/internal/domain/entities"
	"github.com/dinno7/artinux/internal/domain/ports"
	"github.com/dinno7/artinux/internal/domain/services"
	"github.com/dinno7/artinux/pkg/helper"
)

type UploadArtifactUC struct {
	storage       ports.ObjectStorage
	hasher        ports.ChecksumHasher
	fileValidator *services.FileValidator
}

func NewUploadArtifactUC(
	storage ports.ObjectStorage,
	hasher ports.ChecksumHasher,
	fileValidator *services.FileValidator,
) *UploadArtifactUC {
	return &UploadArtifactUC{
		storage:       storage,
		hasher:        hasher,
		fileValidator: fileValidator,
	}
}

func (uc *UploadArtifactUC) Execute(ctx context.Context, filePath string) (string, error) {
	// NOTE: validate file
	validatedFile, err := uc.fileValidator.Validate(filePath)
	if err != nil {
		return "", err
	}

	// NOTE: Create domain entity
	artifact, err := entities.NewArtifact(
		validatedFile.FileName,
		validatedFile.Extension,
		helper.GetRuntimeHostname(),
		helper.GetRuntimeUsername(),
		helper.GetRuntimeOS(),
		helper.GetRuntimeArch(),
		validatedFile.FileSize,
	)
	if err != nil {
		return "", err
	}

	// NOTE: open the file for upload
	file, err := os.Open(validatedFile.FilePath)
	if err != nil {
		return "", domain.ErrInternal.Wrap(err)
	}
	defer file.Close()

	// NOTE: Compute local checksum
	localChecksum, err := uc.hasher.ComputeFromReaderToBase64(file)
	if err != nil {
		return "", err
	}
	artifact.AddChecksum(localChecksum)

	// NOTE: for use again in hasher
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", domain.ErrInternal.Wrap(err)
	}

	// NOTE: Upload to storage
	storageChecksum, err := uc.storage.Upload(ctx, file, artifact)
	if err != nil {
		return "", err
	}

	if storageChecksum != localChecksum {
		return "", domain.ErrFileUploadNotSupportIntegrity
	}

	return artifact.ObjectKey, nil
}
