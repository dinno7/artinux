package artifacts

import (
	"github.com/dinno7/artinux/internal/application/usecases"
	"github.com/dinno7/artinux/pkg/response"
	"github.com/labstack/echo/v4"
)

type UploadArtifactDto struct {
	Arch     string `form:"arch"`
	OS       string `form:"os"`
	Username string `form:"username"`
	Hostname string `form:"hostname"`
}

type UploadedArtifactsDto struct {
	Artifacts []UploadedArtifactDto `json:"artifacts"`
}

type UploadedArtifactDto struct {
	Ok        bool   `json:"ok"`
	FileName  string `json:"file_name"`
	ObjectKey string `json:"object_key,omitempty"`
	Err       string `json:"error,omitempty"`
}

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
//	@Success		201			{object}	response.ResponseSuccessData[UploadedArtifactsDto]
//
//	@Failure		400			{object}	response.ErrorResponseData[any]	"Bad request - missing required field or no files uploaded"
//	@Failure		500			{object}	response.ErrorResponseData[any]	"Internal server error"
//
//	@Router			/artifacts [post]
func (h *ArtifactHTTPHandler) UploadArtifact(c echo.Context) error {
	var payload UploadArtifactDto
	if err := c.Bind(&payload); err != nil {
		return err
	}

	if payload.Arch == "" {
		return response.BadRequestResponse(c, "arch is required")
	}

	if payload.OS == "" {
		err := response.BadRequestResponse(c, "os is required")
		return err
	}

	if payload.Username == "" {
		return response.BadRequestResponse(c, "username is required")
	}

	if payload.Hostname == "" {
		return response.BadRequestResponse(c, "hostname is required")
	}

	incomingFilesHeaders, err := c.MultipartForm()
	if err != nil || incomingFilesHeaders == nil {
		return response.BadRequestResponse(c, "artifact files is required")
	}

	artifactFiles := incomingFilesHeaders.File["artifacts"]
	if len(artifactFiles) == 0 {
		return response.BadRequestResponse(c, "artifacts is required")
	}

	// INFO: Multiple files
	c.Logger().Info("Uploading multiple artifact")
	items := make([]usecases.UploadArtifactItem, len(artifactFiles))
	for i := range artifactFiles {
		// TODO: Sanitize uploaded file name
		art := artifactFiles[i]
		f, err := art.Open()
		if err != nil {
			return err
		}
		defer f.Close()

		items[i] = usecases.UploadArtifactItem{
			FileName:   art.Filename,
			FileSize:   art.Size,
			FileReader: f,
		}
	}
	result, err := h.uploadArtifactsUC.Execute(c.Request().Context(), usecases.UploadArtifactsInput{
		Items:    items,
		OS:       payload.OS,
		Arch:     payload.Arch,
		Hostname: payload.Hostname,
		Username: payload.Username,
	})
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	finalResult := make([]UploadedArtifactsDto, 0, len(artifactFiles))
	for _, uploadedArtifact := range result {
		var err string
		if uploadedArtifact.Err != nil {
			err = uploadedArtifact.Err.Error()
		}
		finalResult = append(finalResult, UploadedArtifactsDto{
			Artifacts: []UploadedArtifactDto{
				{
					Ok:        err == "",
					FileName:  uploadedArtifact.FileName,
					ObjectKey: uploadedArtifact.ObjectKey,
					Err:       err,
				},
			},
		})
	}

	return response.CreatedResponse(c, "Artifact uploaded successfully", finalResult)
}
