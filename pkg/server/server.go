package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/dinno7/artinux/docs"
	"github.com/dinno7/artinux/internal/domain/ports"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	echoSwagger "github.com/swaggo/echo-swagger"
)

const (
	docsPath              = "/docs"
	WriteTimeout          = time.Second * 30
	ReadTimeout           = time.Second * 10
	IdleTimeout           = time.Minute * 1
	ServerShutdownTimeout = time.Second * 5
	GracefulTimeout       = time.Second * 5
)

type Router struct {
	addr        string
	httpHandler *echo.Echo
	logger      ports.Logger
}

// @title			Artinux
// @BasePath		/api/v1
// @description	Artifact Manager
// @termsOfService	http://swagger.io/terms/
//
// @contact.name	Taha Delroba
// @contact.email	tahadlrb7@gmail.com
func NewRouter(addr string, logger ports.Logger) *Router {
	e := echo.New()
	setupMiddlewares(e, logger)

	e.GET(fmt.Sprintf("%s/*", docsPath), echoSwagger.WrapHandler)

	return &Router{
		addr:        addr,
		httpHandler: e,
		logger:      logger,
	}
}

func (r *Router) GetAPIGroup() *echo.Group {
	return r.httpHandler.Group("/api/v1")
}

func (r *Router) ServeHTTP(ctx context.Context) error {
	srv := &http.Server{
		Addr:         r.addr,
		Handler:      r.httpHandler,
		WriteTimeout: WriteTimeout,
		ReadTimeout:  ReadTimeout,
		IdleTimeout:  IdleTimeout,
	}
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			if r.logger != nil {
				r.logger.Error("failed to start server", err)
			}
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), GracefulTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		if r.logger != nil {
			r.logger.Error("failed to stop server", err)
		}
		return err
	}

	return nil
}

func setupMiddlewares(e *echo.Echo, logger ports.Logger) {
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLoggerWithConfig(
		middleware.RequestLoggerConfig{
			LogLatency:       true,
			LogRemoteIP:      true,
			LogHost:          true,
			LogMethod:        true,
			LogURI:           true,
			LogRequestID:     true,
			LogUserAgent:     true,
			LogStatus:        true,
			LogContentLength: true,
			LogResponseSize:  true,
			// forwards error to the global error handler, so it can decide appropriate status code.
			// NB: side-effect of that is - request is now "committed" written to the client. Middlewares up in chain can not
			// change Response status code or response body.
			HandleError: true,
			Skipper: func(c echo.Context) bool {
				return strings.HasPrefix(c.Request().URL.Path, docsPath)
			},
			LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
				if v.Error == nil {
					logger.Info(
						"REQUEST",
						"method", v.Method,
						"uri", v.URI,
						"status", v.Status,
						"latency", v.Latency,
						"host", v.Host,
						"bytes_in", v.ContentLength,
						"bytes_out", v.ResponseSize,
						"user_agent", v.UserAgent,
						"remote_ip", v.RemoteIP,
						"request_id", v.RequestID,
					)
					return nil
				}

				logger.Error(
					"REQUEST_ERROR",
					v.Error,
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
					"latency", v.Latency,
					"host", v.Host,
					"bytes_in", v.ContentLength,
					"bytes_out", v.ResponseSize,
					"user_agent", v.UserAgent,
					"remote_ip", v.RemoteIP,
					"request_id", v.RequestID,
				)
				return nil
			},
		},
	))
	e.Use(middleware.RemoveTrailingSlash())
}
