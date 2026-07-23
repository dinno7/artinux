package artifacts

import (
	"context"

	"github.com/dinno7/artinux/internal/application/usecases"
	_ "github.com/dinno7/artinux/internal/domain/entities"
	"github.com/dinno7/artinux/pkg/response"
	"github.com/labstack/echo/v4"
)

type DeleteArtifactsDto struct {
	ObjectKeys []string `json:"object_keys"`
}

// DeleteArtifacts godoc
//
//	@Summary		Delete artifacts
//	@Description	Deletes one or more artifacts. Supports deletion via object key in the URL path (single delete) or via request body with multiple object keys (batch delete).
//	@Tags			artifacts
//	@Accept			json
//	@Produce		json
//	@Param			object_key	path	string				false	"Single object key to delete (in URL path)"
//	@Param			request		body	DeleteArtifactsDto	false	"Multiple object keys to delete (in request body)"
//	@Success		204			"Artifacts deleted successfully"
//	@Failure		400			{object}	response.ErrorResponseData[any]	"Bad request - invalid payload, empty keys, or deletion error"
//	@Failure		500			{object}	response.ErrorResponseData[any]	"Internal server error"
//	@Router			/artifacts/{object_key} [delete]
//	@Router			/artifacts [delete]
func (h *ArtifactHTTPHandler) DeleteArtifacts(c echo.Context) error {
	ctx := c.Request().Context()

	objectKey, err := extractObjectKeyFromPath(c)
	if err != nil {
		return response.BadRequestResponse(c, "Please provide valid object key")
	}

	if objectKey != "" {
		if err := h.deleteArtifact(ctx, objectKey); err != nil {
			return response.BadRequestResponse(c, err.Error())
		}
		return response.NoContentResponse(c)
	}

	// INFO: Fallback to body with miltiple object keys
	var payload DeleteArtifactsDto
	if err := c.Bind(&payload); err != nil {
		return response.BadRequestResponse(c, "object_keys must be array of keys")
	}

	if len(payload.ObjectKeys) == 0 {
		return response.BadRequestResponse(c, "Please provide object keys")
	}
	if err := h.deleteArtifacts(ctx, payload.ObjectKeys); err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	return response.NoContentResponse(c)
}

func (h *ArtifactHTTPHandler) deleteArtifact(ctx context.Context, objectKey string) error {
	return h.deleteArtifactUC.Execute(ctx, usecases.DeleteArtifactInput{
		ObjectKey: objectKey,
	})
}

func (h *ArtifactHTTPHandler) deleteArtifacts(ctx context.Context, objectKeys []string) error {
	return h.deleteArtifactsUC.Execute(ctx, usecases.DeleteArtifactsInput{
		ObjectKeys: objectKeys,
	})
}
