package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ErrorResponseData represents a standard API error response
//
//	@Description	Standard error response returned by the API
type ErrorResponseData[T any] struct {
	// HTTP status code
	Status int `json:"status"`
	// Error message
	Message string `json:"message"`
	Meta    T      `json:"meta"`
}

// ErrorResponse sends a standardized error response.
func ErrorResponse[T any](c echo.Context, message string, status int, meta ...T) error {
	res := ErrorResponseData[T]{
		Status:  status,
		Message: message,
	}
	if len(meta) > 0 {
		res.Meta = meta[0]
	}
	return c.JSON(status, res)
}

// BadRequestResponse sends a standardized error response with 400 code.
func BadRequestResponse(c echo.Context, message string) error {
	return ErrorResponse[any](c, message, http.StatusBadRequest)
}
