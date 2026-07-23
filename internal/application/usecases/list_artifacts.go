package usecases

import (
	"context"

	"github.com/dinno7/artinux/internal/domain/entities"
	"github.com/dinno7/artinux/internal/domain/ports"
)

type ListArtifactUC struct {
	storage ports.ObjectStorage
	logger  ports.Logger
}

func NewListArtifactUC(
	logger ports.Logger,
	storage ports.ObjectStorage,
) *ListArtifactUC {
	return &ListArtifactUC{
		logger:  logger.With("usecase", "list_artifact"),
		storage: storage,
	}
}

type ListArtifactInput struct {
	Prefix string
	Limit  int
}

func (uc *ListArtifactUC) Execute(
	ctx context.Context,
	input ListArtifactInput,
) ([]*entities.Artifact, error) {
	if input.Limit <= 0 {
		input.Limit = 10
	}

	uc.logger.Info("Getting artifact list", "prefix", input.Prefix, "limit", input.Limit)
	artifacts, err := uc.storage.ListObjects(ctx, input.Prefix, input.Limit)
	if err != nil {
		uc.logger.Error(
			"Failed to get artifacts list",
			err,
			"prefix", input.Prefix,
			"limit", input.Limit,
		)
		return nil, err
	}

	return artifacts, nil
}
