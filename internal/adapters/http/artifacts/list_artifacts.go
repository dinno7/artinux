package artifacts

import (
	"github.com/dinno7/artinux/internal/application/usecases"
	_ "github.com/dinno7/artinux/internal/domain/entities"
	"github.com/dinno7/artinux/pkg/response"
	"github.com/labstack/echo/v5"
)

type ListArtifactDto struct {
	Prefix string `query:"prefix"`
	Limit  int    `query:"limit"`
}

// ListArtifact godoc
//
//	@Summary		List artifacts
//	@Description	Retrieves a filtered, paginated list of artifacts from storage. Supports prefix filtering and result limiting.
//	@Tags			artifacts
//	@Accept			json
//	@Produce		json
//	@Param			prefix	query		string												false	"Filter artifacts by key prefix"	example(artifacts/linux/amd64/)
//	@Param			limit	query		int													false	"Maximum number of results"			example(50)
//	@Success		200		{object}	response.ResponseSuccessData[[]entities.Artifact]	"List retrieved successfully"
//	@Failure		400		{object}	response.ErrorResponseData[any]						"Bad request - invalid parameters"
//	@Failure		500		{object}	response.ErrorResponseData[any]						"Internal server error"
//	@Router			/artifacts [get]
func (h *ArtifactHTTPHandler) ListArtifact(c *echo.Context) error {
	var payload ListArtifactDto

	if err := c.Bind(&payload); err != nil {
		return err
	}

	list, err := h.listArtifactUC.Execute(c.Request().Context(), usecases.ListArtifactInput{
		Prefix: payload.Prefix,
		Limit:  payload.Limit,
	})
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	return response.OkResponse(c, "Fetching list was successfull", list)
}
