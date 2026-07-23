package http

import (
	"context"

	"github.com/dinno7/artinux/pkg/response"
	"github.com/labstack/echo/v4"
)

type Pingable interface {
	Name() string
	Ping(ctx context.Context) error
}

type CommonHTTPHandler struct {
	env       string
	pingables []Pingable
}

func NewCommonHTTPHandler(env string, pingables []Pingable) *CommonHTTPHandler {
	return &CommonHTTPHandler{
		env:       env,
		pingables: pingables,
	}
}

// @Description	Standard API response wrapper with status, message, and data
type dependency struct {
	Name string `json:"name"`
	Ok   bool   `json:"ok"`
}
type healthResult struct {
	Env          string       `json:"env"`
	Dependencies []dependency `json:"dependencies"`
}

//	@Summary		Health check
//	@Description	Runs Ping(ctx) for each configured pingable and reports per-component readiness.
//	@Router			/health [get]
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	response.ResponseSuccessData[healthResult]	"Success"
//
// Health checks all configured pingables and returns their status.
func (h *CommonHTTPHandler) Health(c echo.Context) error {
	dependencies := []dependency{{"HTTP Server", true}}
	for _, pingable := range h.pingables {
		err := pingable.Ping(c.Request().Context())
		dependencies = append(dependencies, dependency{
			Name: pingable.Name(),
			Ok:   err == nil,
		})
	}
	results := &healthResult{
		Env:          h.env,
		Dependencies: dependencies,
	}
	return response.OkResponse(c, "Successfull", results)
}
