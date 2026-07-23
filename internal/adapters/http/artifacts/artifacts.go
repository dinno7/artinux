package artifacts

import "github.com/dinno7/artinux/internal/application/usecases"

type ArtifactHTTPHandler struct {
	uploadArtifactUC   *usecases.UploadArtifactUC
	uploadArtifactsUC  *usecases.UploadArtifactsUC
	listArtifactUC     *usecases.ListArtifactUC
	downloadArtifactUC *usecases.DownloadArtifactUC
	deleteArtifactUC   *usecases.DeleteArtifactUC
	deleteArtifactsUC  *usecases.DeleteArtifactsUC
}

func NewArtifactHTTPHandler(
	uploadArtifactUC *usecases.UploadArtifactUC,
	uploadArtifactsUC *usecases.UploadArtifactsUC,
	listArtifactUC *usecases.ListArtifactUC,
	downloadArtifactUC *usecases.DownloadArtifactUC,
	deleteArtifactUC *usecases.DeleteArtifactUC,
	deleteArtifactsUC *usecases.DeleteArtifactsUC,
) *ArtifactHTTPHandler {
	return &ArtifactHTTPHandler{
		uploadArtifactUC:   uploadArtifactUC,
		uploadArtifactsUC:  uploadArtifactsUC,
		listArtifactUC:     listArtifactUC,
		downloadArtifactUC: downloadArtifactUC,
		deleteArtifactUC:   deleteArtifactUC,
		deleteArtifactsUC:  deleteArtifactsUC,
	}
}
