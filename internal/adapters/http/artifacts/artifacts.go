package artifacts

import "github.com/dinno7/artinux/internal/application/usecases"

type ArtifactHTTPHandler struct {
	uploadArtifactUC   *usecases.UploadArtifactUC
	listArtifactUC     *usecases.ListArtifactUC
	downloadArtifactUC *usecases.DownloadArtifactUC
	deleteArtifactUC   *usecases.DeleteArtifactUC
}

func NewArtifactHTTPHandler(
	uploadArtifactUC *usecases.UploadArtifactUC,
	listArtifactUC *usecases.ListArtifactUC,
	downloadArtifactUC *usecases.DownloadArtifactUC,
	deleteArtifactUC *usecases.DeleteArtifactUC,
) *ArtifactHTTPHandler {
	return &ArtifactHTTPHandler{
		uploadArtifactUC:   uploadArtifactUC,
		listArtifactUC:     listArtifactUC,
		downloadArtifactUC: downloadArtifactUC,
		deleteArtifactUC:   deleteArtifactUC,
	}
}
