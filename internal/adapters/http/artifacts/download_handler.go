package artifacts

import (
	"fmt"
	"io"
	"net/http"

	"github.com/dinno7/artinux/internal/application/usecases"
	"github.com/dinno7/artinux/pkg/response"
	"github.com/labstack/echo/v4"
)

// DownloadArtifact godoc
//
//	@Summary		Download artifact
//	@Description	Dowload a single artifact file with system metadata (arch, os, username, hostname). Supports multipart/form-data with file and form fields.
//	@Tags			artifacts
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			object_key	path		string	true	"Object key"	example(linux/amd64/2026/7/22/071671ad-acbd-4c77-9132-b936ad1187c9_artinux_build.gz)
//
//	@Success		200			{object}	response.ResponseSuccessData[any]
//
//	@Failure		400			{object}	response.ErrorResponseData[any]	"Bad request - missing required field or no files uploaded"
//	@Failure		500			{object}	response.ErrorResponseData[any]	"Internal server error"
//
//	@Router			/artifacts/download/{object_key} [get]
func (h *ArtifactHTTPHandler) DownloadArtifact(c echo.Context) error {
	objectKey, err := extractObjectKeyFromPath(c)
	if err != nil {
		return response.BadRequestResponse(c, "Please provide valid object key")
	}

	output, err := h.downloadArtifactUC.Execute(
		c.Request().Context(),
		usecases.DownloadArtifactInput{
			ObjectKey: objectKey,
		},
	)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	resp := c.Response()

	resp.Header().Set(echo.HeaderContentType, echo.MIMEOctetStream)
	resp.
		Header().
		Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"%s\"", output.Artifact.Name))
	resp.WriteHeader(http.StatusOK)

	if _, err := io.Copy(resp, output.FileReader); err != nil {
		return response.InternalServerResponse(c, err.Error())
	}

	return nil
}
