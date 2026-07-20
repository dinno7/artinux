package artifacts

import (
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"path/filepath"

	"github.com/dinno7/artinux/internal/application/usecases"
	"github.com/dinno7/artinux/pkg/response"
	"github.com/labstack/echo/v4"
)

// UploadArtifact godoc
//
//	@Summary		Upload artifact
//	@Description	Uploads a single artifact file with system metadata (arch, os, username, hostname). Supports multipart/form-data with file and form fields.
//	@Tags			artifacts
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			arch		formData	string	true	"Target architecture (e.g., amd64, arm64)"					example(amd64)
//	@Param			os			formData	string	true	"Target operating system (e.g., linux, darwin, windows)"	example(linux)
//	@Param			username	formData	string	true	"Username of the uploader machine"							example(johndoe)
//	@Param			hostname	formData	string	true	"Hostname of the source machine"							example(web-server-01)
//	@Param			artifacts	formData	[]file	true	"Artifact file to upload (multipart)"
//	@Success		201			{object}	response.ResponseSuccessData[map[string]string]
//
//	@Failure		400			{object}	response.ErrorResponseData[any]	"Bad request - missing required field or no files uploaded"
//	@Failure		500			{object}	response.ErrorResponseData[any]	"Internal server error"
//
//	@Router			/artifacts [post]
func (h *ArtifactHTTPHandler) UploadArtifact(c echo.Context) error {
	arch := c.FormValue("arch")
	osName := c.FormValue("os")
	username := c.FormValue("username")
	hostname := c.FormValue("hostname")

	if arch == "" {
		return response.BadRequestResponse(c, "arch is required")
	}

	if osName == "" {
		err := response.BadRequestResponse(c, "os is required")
		return err
	}

	if username == "" {
		return response.BadRequestResponse(c, "username is required")
	}

	if hostname == "" {
		return response.BadRequestResponse(c, "hostname is required")
	}

	incomingFilesHeaders, err := c.MultipartForm()
	if err != nil {
		return err
	}

	artifactFiles := incomingFilesHeaders.File["artifacts"]
	if incomingFilesHeaders == nil || len(artifactFiles) == 0 {
		return response.BadRequestResponse(c, "artifacts is required")
	}

	requestedFiles := []string{}

	tmpDirPath := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("artinux.%d", rand.Int32()),
	)
	if err := os.MkdirAll(tmpDirPath, 0o700); err != nil {
		return err
	}
	defer os.RemoveAll(tmpDirPath)

	for i := range artifactFiles {
		artifactFile := artifactFiles[i]
		incomingFile, err := artifactFile.Open()
		if err != nil {
			return err
		}

		// TODO: Sanitize uploaded file name
		tmpFilePath := filepath.Join(tmpDirPath, artifactFile.Filename)
		tmpFile, err := os.OpenFile(
			tmpFilePath,
			os.O_CREATE|os.O_WRONLY,
			0o600,
		)
		if err != nil {
			return err
		}

		if _, err := io.Copy(tmpFile, incomingFile); err != nil {
			c.Logger().Error("failed store incomming file", err)
			continue
		}

		requestedFiles = append(requestedFiles, tmpFile.Name())
		incomingFile.Close()
		tmpFile.Close()
	}

	if len(requestedFiles) == 1 {
		objKey, err := h.uploadArtifactUC.Execute(
			c.Request().Context(),
			usecases.UploadArtifactInput{
				FilePath: requestedFiles[0],
				OS:       osName,
				Arch:     arch,
				Hostname: hostname,
				Username: username,
			},
		)
		if err != nil {
			return err
		}
		return response.CreatedResponse(c, "Artifact uploaded successfully", map[string]string{
			"object_key": objKey,
		})
	}

	// TODO: Batch upload

	return echo.ErrInternalServerError
}
