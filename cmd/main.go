package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	httpCommon "github.com/dinno7/artinux/internal/adapters/http"
	httpArtifacts "github.com/dinno7/artinux/internal/adapters/http/artifacts"
	"github.com/dinno7/artinux/internal/application/usecases"
	"github.com/dinno7/artinux/internal/domain/services"
	hasher "github.com/dinno7/artinux/internal/infrastructure/checksum_hasher"
	"github.com/dinno7/artinux/internal/infrastructure/config"
	"github.com/dinno7/artinux/internal/infrastructure/logger"
	"github.com/dinno7/artinux/internal/infrastructure/storage"
	"github.com/dinno7/artinux/pkg/response"
	"github.com/dinno7/artinux/pkg/server"
	"github.com/labstack/echo/v4"
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
		logger.Fatal(
			"failed to connect MinIO",
			err,
		)
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
	uploadArtifactsUC := usecases.NewUploadArtifactsUC(
		logger,
		objStorage,
		checksumHasher,
		fileValidator,
	)
	downloadArtifactUC := usecases.NewDownloadArtifactUC(logger, objStorage, checksumHasher)
	deleteArtifactUC := usecases.NewDeleteArtifactUC(logger, objStorage)
	deleteArtifactsUC := usecases.NewDeleteArtifactsUC(logger, objStorage)
	listArtifactUC := usecases.NewListArtifactUC(logger, objStorage)

	commonHTTPHandler := httpCommon.NewCommonHTTPHandler(cfg.Env, []httpCommon.Pingable{
		objStorage,
	})
	artifactHTTPHandler := httpArtifacts.NewArtifactHTTPHandler(
		uploadArtifactUC,
		uploadArtifactsUC,
		listArtifactUC,
		downloadArtifactUC,
		deleteArtifactUC,
		deleteArtifactsUC,
	)

	addr := fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port)
	router := server.NewRouter(addr, logger)
	apiGroup := router.GetAPIGroup()

	apiGroup.GET("/health", commonHTTPHandler.Health)
	apiGroup.DELETE("/clear", func(c echo.Context) error {
		if err := objStorage.ClearBucket(c.Request().Context()); err != nil {
			logger.Error("Failed clear bucket", err)
			return response.InternalServerResponse(c, "failed clear bucket")
		}
		return nil
	})

	artifactsGroup := apiGroup.Group("/artifacts")
	artifactsGroup.GET("", artifactHTTPHandler.ListArtifact)
	artifactsGroup.GET("/download/*", artifactHTTPHandler.DownloadArtifact)
	artifactsGroup.POST("", artifactHTTPHandler.UploadArtifact)
	artifactsGroup.DELETE("*", artifactHTTPHandler.DeleteArtifacts)

	logger.Info("Server is running", "addr", addr)
	if err := router.ServeHTTP(ctx); err != nil {
		logger.Fatal("failed start http server", err)
	}
	<-ctx.Done()
}
