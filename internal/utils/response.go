package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error string      `json:"error"`
	Code  string      `json:"code,omitempty"`
	Info  interface{} `json:"info,omitempty"`
}

func JSONError(c *gin.Context, status int, message string) {
	c.JSON(status, ErrorResponse{Error: message})
}

func JSONErrorWithCode(c *gin.Context, status int, code, message string) {
	c.JSON(status, ErrorResponse{Error: message, Code: code})
}

func JSONValidationError(c *gin.Context, err error) {
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		details := make([]string, 0, len(validationErrs))
		for _, fieldErr := range validationErrs {
			details = append(details, formatValidationError(fieldErr))
		}
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "validation failed",
			Code:  "VALIDATION_ERROR",
			Info:  details,
		})
		return
	}

	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error: "invalid request payload",
		Code:  "BAD_REQUEST",
		Info:  err.Error(),
	})
}

func formatValidationError(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fieldErr.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", fieldErr.Field())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", fieldErr.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fieldErr.Field(), fieldErr.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", fieldErr.Field(), fieldErr.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of [%s]", fieldErr.Field(), fieldErr.Param())
	default:
		return fmt.Sprintf("%s is invalid", fieldErr.Field())
	}
}
