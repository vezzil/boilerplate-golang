package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Success writes a standard success JSON payload with 200 status code
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// SuccessWithStatus writes a success JSON payload with custom status code
func SuccessWithStatus(c *gin.Context, statusCode int, data any, message string) {
	if message == "" {
		message = "Operation completed successfully"
	}

	c.JSON(statusCode, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// Error writes a standard error JSON payload
func Error(c *gin.Context, code int, message string) {
	if code <= 0 {
		code = http.StatusInternalServerError
	}

	c.JSON(code, gin.H{
		"success": false,
		"error":   message,
	})
}

// ValidationError writes a validation error JSON payload
func ValidationError(c *gin.Context, errors map[string]string) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{
		"success": false,
		"error":   "Validation failed",
		"errors":  errors,
	})
}

// NotFound writes a 404 not found error
func NotFound(c *gin.Context, resource string) {
	if resource == "" {
		resource = "Resource"
	}
	Error(c, http.StatusNotFound, resource+" not found")
}

// Unauthorized writes a 401 unauthorized error
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	Error(c, http.StatusUnauthorized, message)
}

// Forbidden writes a 403 forbidden error
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "Forbidden"
	}
	Error(c, http.StatusForbidden, message)
}

// BadRequest writes a 400 bad request error
func BadRequest(c *gin.Context, message string) {
	if message == "" {
		message = "Bad request"
	}
	Error(c, http.StatusBadRequest, message)
}

// InternalServerError writes a 500 internal server error
func InternalServerError(c *gin.Context, message string) {
	if message == "" {
		message = "Internal server error"
	}
	Error(c, http.StatusInternalServerError, message)
}

// Created writes a 201 created response
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    data,
	})
}

// NoContent writes a 204 no content response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
