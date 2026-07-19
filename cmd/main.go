package main

import (
	"context"
	"os"
	"os/signal"

	httpCommon "github.com/dinno7/artinux/internal/adapters/http"
	httpArtifacts "github.com/dinno7/artinux/internal/adapters/http/artifacts"
	"github.com/dinno7/artinux/internal/application/usecases"
	"github.com/dinno7/artinux/internal/domain"
	"github.com/dinno7/artinux/internal/domain/services"
	hasher "github.com/dinno7/artinux/internal/infrastructure/checksum_hasher"
	"github.com/dinno7/artinux/internal/infrastructure/config"
	"github.com/dinno7/artinux/internal/infrastructure/logger"
	"github.com/dinno7/artinux/internal/infrastructure/storage"
)

func main() {
	ctx, cancle := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancle()

	cfg, err := config.Get()
	if err != nil {
		panic(err)
	}

	logger := logger.NewLogger(logger.LoggerConfig{
		Level:  cfg.Logging.Level,
		Format: cfg.Logging.Format,
	})

	objStorage, err := storage.NewMinIOStorage(
		storage.MinIOConfig{
			Region:           cfg.ObjectStorage.Region,
			UseSSL:           cfg.ObjectStorage.UseSSL,
			Endpoint:         cfg.ObjectStorage.Endpoint,
			BucketName:       cfg.ObjectStorage.BucketName,
			AccessKeyID:      cfg.ObjectStorage.Username,
			HealthInterval:   cfg.ObjectStorage.HealthInterval,
			SecretAccessKey:  cfg.ObjectStorage.Password,
			MaxUploadRetries: cfg.ObjectStorage.MaxUploadRetries,
		},
	)
	if err != nil {
		logger.Fatal("failed to connect MinIO", domain.ErrStorageUnavailable.Wrap(err))
	}

	checksumHasher := hasher.NewSha256Hasher()
	fileValidator := services.NewFileValidator(
		cfg.Upload.AllowedFileExtensions,
		cfg.Upload.MaxSizeMB,
	)

	uploadArtifactUC := usecases.NewUploadArtifactUC(
		logger,
		objStorage,
		checksumHasher,
		fileValidator,
	)
	downloadArtifactUC := usecases.NewDownloadArtifactUC(logger, objStorage, checksumHasher)
	deleteArtifactUC := usecases.NewDeleteArtifactUC(logger, objStorage)
	listArtifactUC := usecases.NewListArtifactUC(logger, objStorage)

	commonHTTPHandler := httpCommon.NewCommonHTTPHandler([]httpCommon.Pingable{
		objStorage,
	})
	artifactHTTPHandler := httpArtifacts.NewArtifactHTTPHandler(
		uploadArtifactUC,
		listArtifactUC,
		downloadArtifactUC,
		deleteArtifactUC,
	)

	router := NewRouter("localhost:7000", logger)
	apiGroup := router.GetAPIGroup()

	apiGroup.GET("/health", commonHTTPHandler.Health)
	apiGroup.GET("/artifacts", artifactHTTPHandler.ListArtifact)
	apiGroup.POST("/artifacts", artifactHTTPHandler.UploadArtifact)
	apiGroup.DELETE("/artifacts", artifactHTTPHandler.DeleteArtifacts)

	logger.Info("Server is running on port 7000")
	if err := router.ServeHTTP(ctx); err != nil {
		logger.Fatal("failed start http server", err)
	}
	<-ctx.Done()
}
