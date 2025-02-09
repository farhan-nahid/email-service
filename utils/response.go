package utils

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Success writes a JSON response with the given status code, data, and message
func Success(c *gin.Context, status int, data interface{}, message string) {
	c.JSON(status, gin.H{
		"data":    data,
		"success": true,
		"message": message,
	})
	c.Abort()
}


// Error writes an error response with the given status code and error message
func Error(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{
		"success": false,
		"message": "Error",
		"error": err.Error(),
	})
	c.Abort()
}


// ValidatorError writes a validation error response with the given status code and validation errors
func ValidatorError(c *gin.Context, status int, errors validator.ValidationErrors) {
	var errorMessage []string

	for _, err := range errors {
		fmt.Println(err)
		switch err.ActualTag() {
		case "required":
			errorMessage = append(errorMessage, err.Field() + " is required")
		case "email_address":
			errorMessage = append(errorMessage, err.Field() + " is not a valid email")
		case "status":
			errorMessage = append(errorMessage, err.Field() + " is not a valid status")
		case "source":
			errorMessage = append(errorMessage, err.Field() + " is not a valid source")
		case "website":
			errorMessage = append(errorMessage, err.Field() + " is not a valid website")
		case "uuid":
			errorMessage = append(errorMessage, err.Field() + " is not a valid UUID")
		default:
			errorMessage = append(errorMessage, err.Field() + " is not valid")
		}
	}

	c.JSON(status, gin.H{
		"success": false,
		"message": "Validation error",
		"errors":  strings.Join(errorMessage, ", "),
	})
	c.Abort()
}
