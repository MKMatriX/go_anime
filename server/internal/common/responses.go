package common

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type ApiResponse map[string]any

type ValidationError struct {
	Error     string `json:error`
	Field     string `json: field`
	Condition string `json: condition`
}

type JSONSuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type JSONFailedValidationResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Errors  []*ValidationError `json:"errors"`
}

type JSONErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func SendSuccessResponse(c *echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusOK, JSONSuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SendFailedValidateResponse(c *echo.Context, errors []*ValidationError) error {
	return c.JSON(http.StatusBadRequest, JSONFailedValidationResponse{
		Success: false,
		Message: "Input data validation failed",
		Errors:  errors,
	})
}

func SendErrorResponse(c *echo.Context, message string, statusCode int) error {
	return c.JSON(statusCode, JSONErrorResponse{
		Success: false,
		Message: message,
	})
}

func SendBadRequestResponse(c *echo.Context, message string) error {
	return SendErrorResponse(c, message, http.StatusBadRequest)
}

func SendForbiddenResponse(c *echo.Context) error {
	return SendErrorResponse(c, "Access forbidden", http.StatusForbidden)
}

func SendUnauthorizedResponse(c *echo.Context, message string) error {
	return SendErrorResponse(c, message, http.StatusUnauthorized)
}

func SendNotFoundResponse(c *echo.Context, message string) error {
	if message == "" {
		message = "Entity not found"
	}
	return SendErrorResponse(c, message, http.StatusNotFound)
}
