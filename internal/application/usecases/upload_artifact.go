package usecases

import (
	"context"
	"io"
	"os"

	"github.com/dinno7/artinux/internal/domain"
	"github.com/dinno7/artinux/internal/domain/entities"
	"github.com/dinno7/artinux/internal/domain/ports"
	"github.com/dinno7/artinux/internal/domain/services"
)

type UploadArtifactUC struct {
	storage       ports.ObjectStorage
	hasher        ports.ChecksumHasher
	logger        ports.Logger
	fileValidator *services.FileValidator
}

func NewUploadArtifactUC(
	logger ports.Logger,
	storage ports.ObjectStorage,
	hasher ports.ChecksumHasher,
	fileValidator *services.FileValidator,
) *UploadArtifactUC {
	return &UploadArtifactUC{
		storage:       storage,
		logger:        logger.With("usecase", "upload_artifact"),
		hasher:        hasher,
		fileValidator: fileValidator,
	}
}

type UploadArtifactInput struct {
	FilePath string
	Hostname string
	Username string
	OS       string
	Arch     string
}

func (uc *UploadArtifactUC) Execute(
	ctx context.Context,
	input UploadArtifactInput,
) (string, error) {
	// NOTE: validate file
	uc.logger.Info("Validate incomming file", "path", input.FilePath)
	validatedFile, err := uc.fileValidator.Validate(input.FilePath)
	if err != nil {
		uc.logger.Error("File validation failed", err, "path", input.FilePath)
		return "", err
	}

	makeEmptyUnknown(&input.Hostname)
	makeEmptyUnknown(&input.Username)
	makeEmptyUnknown(&input.OS)
	makeEmptyUnknown(&input.Arch)

	// NOTE: Create domain entity
	uc.logger.Info("Creating domain entity")
	artifact, err := entities.NewArtifact(
		validatedFile.FileName,
		validatedFile.Extension,
		input.Hostname,
		input.Username,
		input.OS,
		input.Arch,
		validatedFile.FileSize,
	)
	if err != nil {
		uc.logger.Error("Failed to craeting domain enity", err, "path", input.FilePath)
		return "", err
	}

	// NOTE: open the file for upload
	uc.logger.Info("Opening file to upload it", "path", input.FilePath)
	file, err := os.Open(validatedFile.FilePath)
	if err != nil {
		uc.logger.Error("Failed to opening file for upload", err, "path", input.FilePath)
		return "", domain.ErrInternal.Wrap(err)
	}
	defer file.Close()

	// NOTE: Compute local checksum
	uc.logger.Info("Computing file checksum", "path", input.FilePath)
	localChecksum, err := uc.hasher.ComputeFromReaderToBase64(file)
	if err != nil {
		uc.logger.Error("Failed to compute checksum", err, "path", input.FilePath)
		return "", err
	}
	uc.logger.Info("Checksum computed", "path", input.FilePath, "checksum", localChecksum)
	artifact.AddChecksum(localChecksum)

	// NOTE: for use again in hasher
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		uc.logger.Error("Failed to reset opened file reader's pointer", err, "path", input.FilePath)
		return "", domain.ErrInternal.Wrap(err)
	}

	// NOTE: Upload to storage
	uc.logger.Info("Uploading file", "path", input.FilePath, "object_key", artifact.ObjectKey)
	storageChecksum, err := uc.storage.Upload(ctx, file, artifact)
	if err != nil {
		uc.logger.Error(
			"Failed to uploading file",
			err,
			"path", input.FilePath,
			"object_key", artifact.ObjectKey,
		)
		return "", err
	}

	if storageChecksum != localChecksum {
		uc.logger.Error(
			"Uploaded & Local checksums are not same",
			nil,
			"object_key", artifact.ObjectKey,
			"uploaded_checksum", storageChecksum,
			"local_checksum", localChecksum,
		)
		return "", domain.ErrFileUploadNotSupportIntegrity
	}

	return artifact.ObjectKey, nil
}
