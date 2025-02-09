package middleware

import (
	"net/http"

	"github.com/farhan-nahid/email-service/models"
	"github.com/farhan-nahid/email-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BindAndValidate is a middleware that binds JSON from request body to struct and validates the struct
func BindAndValidate[T interface{}]() gin.HandlerFunc {
	var validate = validator.New() // Initialize validator instance
	models.RegisterCustomValidations(validate) // Register custom validations

	return func(c *gin.Context) {
		var input T
		// Bind JSON from request body to struct
		if err := c.ShouldBindJSON(&input); err != nil {	
			utils.Error(c, http.StatusBadRequest, err)
			return
		}

		// Validate struct
		if err := validate.Struct(input); err != nil {
			utils.ValidatorError(c, http.StatusBadRequest, err.(validator.ValidationErrors))
			return
		}

		// Pass validated data to the next handler
		c.Set("validatedData", input)
		c.Next()
	}
}