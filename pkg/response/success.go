package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ResponseSuccessData[T any] struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// SuccessResponse sends a standardized success response
func SuccessResponse[T any](c echo.Context, message string, data T, status int) error {
	return c.JSON(status, ResponseSuccessData[T]{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

// OkResponse sends a standardized success response with 200 status
func OkResponse(c echo.Context, message string, data any) error {
	return SuccessResponse(c, message, data, http.StatusOK)
}

// CreatedResponse sends a standardized success response with 201 status
func CreatedResponse(c echo.Context, message string, data any) error {
	return SuccessResponse(c, message, data, http.StatusCreated)
}

// NoContentResponse sends a standardized success response with 204 status
func NoContentResponse(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
