package artifacts

import (
	"github.com/dinno7/artinux/internal/application/usecases"
	_ "github.com/dinno7/artinux/internal/domain/entities"
	"github.com/dinno7/artinux/pkg/response"
	"github.com/labstack/echo/v5"
)

type DeleteArtifactsDto struct {
	ObjectKeys []string `json:"object_keys"`
}

// DeleteArtifacts godoc
//
//	@Summary		Delete artifacts
//	@Description	Deletes one or more artifacts by their object keys. Currently single-key deletion only; batch delete is not yet implemented.
//	@Tags			artifacts
//	@Accept			json
//	@Produce		json
//	@Param			request	body	DeleteArtifactsDto	true	"Object keys to delete"
//	@Success		204		"Artifact deleted successfully"
//	@Failure		400		{object}	response.ErrorResponseData[any]	"Bad request - invalid payload, empty keys, or batch delete not implemented"
//	@Failure		500		{object}	response.ErrorResponseData[any]	"Internal server error"
//	@Router			/artifacts [delete]
func (h *ArtifactHTTPHandler) DeleteArtifacts(c *echo.Context) error {
	var payload DeleteArtifactsDto
	if err := c.Bind(&payload); err != nil {
		return response.BadRequestResponse(c, "object_keys must be array of keys")
	}

	if len(payload.ObjectKeys) > 1 {
		return response.BadRequestResponse(c, "Batch delete not implemented yet")
	}

	if len(payload.ObjectKeys) == 1 {
		err := h.deleteArtifactUC.Execute(c.Request().Context(), usecases.DeleteArtifactInput{
			ObjectKey: payload.ObjectKeys[0],
		})
		if err != nil {
			return response.BadRequestResponse(c, err.Error())
		}
		return response.NoContentResponse(c)
	}

	// TODO: Add batch delete

	return nil
}
