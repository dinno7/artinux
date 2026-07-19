package usecases

import (
	"context"

	"github.com/dinno7/artinux/internal/domain/ports"
)

type DeleteArtifactUC struct {
	storage ports.ObjectStorage
	logger  ports.Logger
}

func NewDeleteArtifactUC(
	logger ports.Logger,
	storage ports.ObjectStorage,
) *DeleteArtifactUC {
	return &DeleteArtifactUC{
		logger:  logger.With("usecase", "delete_artifact"),
		storage: storage,
	}
}

type DeleteArtifactInput struct {
	ObjectKey string
}

func (uc *DeleteArtifactUC) Execute(
	ctx context.Context,
	input DeleteArtifactInput,
) error {
	uc.logger.Info("Deleting artifact", "object_key", input.ObjectKey)
	err := uc.storage.DeleteObject(ctx, input.ObjectKey)
	if err != nil {
		uc.logger.Error(
			"Failed to delete artifact",
			err,
			"object_key", input.ObjectKey,
		)
		return err
	}
	return nil
}
