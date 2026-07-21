package usecases

import (
	"context"
	"io"

	"github.com/dinno7/artinux/internal/domain"
	"github.com/dinno7/artinux/internal/domain/entities"
	"github.com/dinno7/artinux/internal/domain/ports"
	"github.com/dinno7/artinux/internal/domain/services"
	"github.com/dinno7/artinux/pkg/helper"
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
	makeEmptyUnknown(&input.Hostname)
	makeEmptyUnknown(&input.Username)
	makeEmptyUnknown(&input.OS)
	makeEmptyUnknown(&input.Arch)

	if !helper.IsValidOSAndArch(input.OS, input.Arch) {
		return "", domain.ErrInvalidOsOrArch
	}

	validatedFile, err := uc.fileValidator.Validate(input.FilePath)
	if err != nil {
		uc.logger.Error("File validation failed", err, "path", input.FilePath)
		return "", err
	}

	file := validatedFile.File
	defer file.Close()

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

	// NOTE: Compute local checksum
	uc.logger.Info("Computing file checksum", "path", input.FilePath)
	localChecksum, err := uc.hasher.ComputeFromReaderToBase64(file)
	if err != nil {
		uc.logger.Error("Failed to compute checksum", err, "path", input.FilePath)
		return "", err
	}
	uc.logger.Info("Checksum computed", "path", input.FilePath, "checksum", localChecksum)
	artifact.AddChecksum(localChecksum)

	// NOTE: for use again
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
		err := uc.storage.DeleteObject(ctx, artifact.ObjectKey)
		if err != nil {
			uc.logger.Warn(
				"Failed to delete object from storage which uploaded via invalid checksum",
				nil,
				"object_key", artifact.ObjectKey,
				"uploaded_checksum", storageChecksum,
				"local_checksum", localChecksum,
			)
		}
		return "", domain.ErrFileUploadNotSupportIntegrity
	}

	return artifact.ObjectKey, nil
}
