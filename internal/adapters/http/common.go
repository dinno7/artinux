package http

import (
	"context"

	"github.com/dinno7/artinux/pkg/response"
	"github.com/labstack/echo/v5"
)

type Pingable interface {
	Name() string
	Ping(ctx context.Context) error
}

type CommonHTTPHandler struct {
	pingables []Pingable
}

func NewCommonHTTPHandler(pingables []Pingable) *CommonHTTPHandler {
	return &CommonHTTPHandler{
		pingables: pingables,
	}
}

// @Description	Standard API response wrapper with status, message, and data
type healthResult struct {
	Name string `json:"name"`
	Ok   bool   `json:"ok"`
}

//	@Summary		Health check
//	@Description	Runs Ping(ctx) for each configured pingable and reports per-component readiness.
//	@Router			/api/v1/health [get]
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	response.ResponseSuccessData[healthResult]	"Success"
//
// Health checks all configured pingables and returns their status.
func (h *CommonHTTPHandler) Health(c *echo.Context) error {
	results := []*healthResult{}
	for _, pingable := range h.pingables {
		err := pingable.Ping(c.Request().Context())
		results = append(results, &healthResult{
			Name: pingable.Name(),
			Ok:   err == nil,
		})
	}
	return response.OkResponse(c, "Successfull", results)
}
