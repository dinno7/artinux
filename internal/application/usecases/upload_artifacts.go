package usecases

import (
	"context"
	"io"
	"strings"
	"sync"

	"github.com/dinno7/artinux/internal/domain"
	"github.com/dinno7/artinux/internal/domain/entities"
	"github.com/dinno7/artinux/internal/domain/ports"
	"github.com/dinno7/artinux/internal/domain/services"
	"github.com/dinno7/artinux/pkg/helper"
)

type UploadArtifactsUC struct {
	storage       ports.ObjectStorage
	hasher        ports.ChecksumHasher
	logger        ports.Logger
	fileValidator *services.FileValidator
}

func NewUploadArtifactsUC(
	logger ports.Logger,
	storage ports.ObjectStorage,
	hasher ports.ChecksumHasher,
	fileValidator *services.FileValidator,
) *UploadArtifactsUC {
	return &UploadArtifactsUC{
		storage:       storage,
		hasher:        hasher,
		logger:        logger.With("usecase", "upload_artifacts"),
		fileValidator: fileValidator,
	}
}

type UploadArtifactItem struct {
	// FilePath string
	FileName   string
	FileSize   int64
	FileReader io.ReadSeekCloser
}

type UploadArtifactsInput struct {
	Hostname string
	Username string
	OS       string
	Arch     string
	Items    []UploadArtifactItem
}

type UploadArtifactsResult struct {
	FileName  string
	ObjectKey string
	Err       error
}

func (uc *UploadArtifactsUC) Execute(
	ctx context.Context,
	input UploadArtifactsInput,
) ([]UploadArtifactsResult, error) {
	if len(input.Items) == 0 {
		return nil, domain.ErrStorageEmptyArtifactToUpload
	}

	input.OS = strings.TrimSpace(strings.ToLower(input.OS))
	input.Arch = strings.TrimSpace(strings.ToLower(input.Arch))

	if !helper.IsValidOSAndArch(input.OS, input.Arch) {
		err := domain.ErrInvalidOsOrArch
		uc.logger.Error(
			"OS or Arch validation failed",
			err,
			"os", input.OS,
			"arch", input.Arch,
		)
		return nil, domain.ErrInvalidOsOrArch
	}

	makeEmptyUnknown(&input.Hostname)
	makeEmptyUnknown(&input.Username)

	type preparedArtifact struct {
		artifact      *entities.Artifact
		file          io.ReadSeekCloser
		localChecksum string
		resultIndex   int // tracks original order for result mapping
	}

	results := make([]UploadArtifactsResult, len(input.Items))
	preparedUploads := make([]preparedArtifact, 0, len(input.Items))
	objectKeysToDelete := []string{} // INFO: Objects uploaded but not valid(uploaded checksum != local one) and must delete

	// INFO: Clean up all opened files when done
	defer func() {
		for _, p := range preparedUploads {
			p.file.Close()
		}

		// INFO: Delteting objects with invalid checksum
		if len(objectKeysToDelete) > 0 {
			if err := uc.storage.DeleteBatch(ctx, objectKeysToDelete); err != nil {
				uc.logger.Warn(
					"Failed to delete objects with invalid checksum",
					"object_keys",
					objectKeysToDelete,
				)
			}
		}
	}()

	for i, item := range input.Items {
		// NOTE: Validate file
		uc.logger.Info("Validate incoming file", "file_name", item.FileName)

		ext, err := uc.fileValidator.ValidateAndGetExt(item.FileName, item.FileSize)
		if err != nil {
			uc.logger.Error("File validation failed", err, "file_name", item.FileName)
			results[i] = UploadArtifactsResult{Err: err, FileName: item.FileName}
			continue
		}

		file := item.FileReader

		// NOTE: Create domain entity
		artifact, err := entities.NewArtifact(
			item.FileName,
			ext,
			input.Hostname,
			input.Username,
			input.OS,
			input.Arch,
			item.FileSize,
		)
		if err != nil {
			uc.logger.Error("Failed to create domain entity", err, "file_name", item.FileName)
			file.Close()
			results[i] = UploadArtifactsResult{Err: err, FileName: item.FileName}
			continue
		}

		// INFO: Compute local checksum
		localChecksum, err := uc.hasher.ComputeFromReaderToBase64(file)
		if err != nil {
			uc.logger.Error("Failed to compute checksum", err, "file_name", item.FileName)
			file.Close()
			results[i] = UploadArtifactsResult{Err: err, FileName: item.FileName}
			continue
		}

		// INFO: Reset reader for upload
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			uc.logger.Error("Failed to reset file reader", err, "file_name", item.FileName)
			file.Close()
			results[i] = UploadArtifactsResult{
				Err:      domain.ErrInternal.Wrap(err),
				FileName: item.FileName,
			}
			continue
		}

		artifact.AddChecksum(localChecksum)

		preparedUploads = append(preparedUploads, preparedArtifact{
			artifact:      artifact,
			file:          file,
			localChecksum: localChecksum,
			resultIndex:   i,
		})

		results[i] = UploadArtifactsResult{ObjectKey: artifact.ObjectKey, FileName: item.FileName}
	}

	// INFO: If nothing passed validation, return early
	if len(preparedUploads) == 0 {
		return results, nil
	}

	uploadWG := new(sync.WaitGroup)
	uploadedChannel := make(chan struct {
		preparedIndex int
		err           error
	})
	uc.logger.Info("Uploading objects", "count", len(preparedUploads))
	for i, p := range preparedUploads {
		uploadWG.Add(1)
		go func() {
			defer uploadWG.Done()

			uc.logger.Info("Uploading object", "object_key", p.artifact.ObjectKey)
			uploadedChecksum, err := uc.storage.Upload(ctx, p.file, p.artifact)
			if err != nil {
				uc.logger.Error(
					"Failed to uploading file",
					err,
					"file_name", p.artifact.Name,
					"object_key", p.artifact.ObjectKey,
				)
				uploadedChannel <- struct {
					preparedIndex int
					err           error
				}{
					preparedIndex: i,
					err:           err,
				}
				return
			}

			// NOTE: Verify integrity
			uc.logger.Info(
				"Upload completed, verifying integrity",
				"object_key", p.artifact.ObjectKey,
			)
			if uploadedChecksum != p.localChecksum {
				err := domain.ErrFileChecksumNotSame
				uc.logger.Error(
					"Checksum mismatch after upload",
					err,
					"object_key", p.artifact.ObjectKey,
					"local_checksum", p.localChecksum,
					"uploaded_checksum", uploadedChecksum,
				)
				uploadedChannel <- struct {
					preparedIndex int
					err           error
				}{
					preparedIndex: i,
					err:           err,
				}
				return
			}

			uploadedChannel <- struct {
				preparedIndex int
				err           error
			}{
				preparedIndex: i,
				err:           nil,
			}
		}()
	}

	go func() {
		uploadWG.Wait()
		close(uploadedChannel)
	}()

	for res := range uploadedChannel {
		uploadedData := preparedUploads[res.preparedIndex]
		objectKey := uploadedData.artifact.ObjectKey

		results[uploadedData.resultIndex].ObjectKey = objectKey
		results[uploadedData.resultIndex].Err = res.err

		if res.err != nil {
			objectKeysToDelete = append(objectKeysToDelete, objectKey)
		}
	}

	return results, nil
}
