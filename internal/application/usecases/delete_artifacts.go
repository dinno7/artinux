package usecases

import (
	"context"

	"github.com/dinno7/artinux/internal/domain/ports"
)

type DeleteArtifactsUC struct {
	storage ports.ObjectStorage
	logger  ports.Logger
}

func NewDeleteArtifactsUC(
	logger ports.Logger,
	storage ports.ObjectStorage,
) *DeleteArtifactsUC {
	return &DeleteArtifactsUC{
		logger:  logger.With("usecase", "delete_artifact"),
		storage: storage,
	}
}

type DeleteArtifactsInput struct {
	ObjectKeys []string
}

func (uc *DeleteArtifactsUC) Execute(
	ctx context.Context,
	input DeleteArtifactsInput,
) error {
	uc.logger.Info("Deleting artifacts", "object_keys", input.ObjectKeys)

	err := uc.storage.DeleteBatch(ctx, input.ObjectKeys)
	if err != nil {
		uc.logger.Error(
			"Failed to delete artifacts",
			err,
			"object_keys", input.ObjectKeys,
		)
		return err
	}
	return nil
}
