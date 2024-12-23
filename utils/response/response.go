package response

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)


func WriteJSON(c *gin.Context, status int, data interface{})  {
	c.JSON(status, gin.H{
		"success": true,
		"message": "Success",
		"data": data,
	})
}


func WriteError(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{
		"success": false,
		"message": "Error",
		"error": err.Error(),
	})
	c.Abort()
	
}

func ValidatorError(c *gin.Context, status int, errors validator.ValidationErrors) {
	var errorMessage []string

	for _, err := range errors {
		fmt.Println(err)
		switch err.ActualTag() {
		case "required":
			errorMessage = append(errorMessage, err.Field() + " is required")
		case "email":
			errorMessage = append(errorMessage, err.Field() + " is not a valid email")
		case "status":
			errorMessage = append(errorMessage, err.Field() + " is not a valid status")
		case "source":
			errorMessage = append(errorMessage, err.Field() + " is not a valid source")
		case "website":
			errorMessage = append(errorMessage, err.Field() + " is not a valid website")
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
