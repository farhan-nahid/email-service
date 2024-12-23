package controllers

import (
	"errors"
	"net/http"

	"github.com/farhan-nahid/email-service/initializers"
	"github.com/farhan-nahid/email-service/models"
	"github.com/farhan-nahid/email-service/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateEmail(c *gin.Context) {
	// Retrieve the validated email data from the context
	validatedData, exists := c.Get("validatedData")
	if !exists {
		// If no validated data exists in the context, return an error response
		response.WriteError(c, http.StatusBadRequest, errors.New("data not found in context"))
		return
	}

	// Type assert validatedData into the Email model
	emailData, ok := validatedData.(models.Email)
	if !ok {
		// If the type assertion fails, return an error response
		response.WriteError(c, http.StatusBadRequest, errors.New("invalid data format"))
		return
	}

	// Create a new email instance using the validated data
	newEmail := models.Email{
		CompanyUUID: emailData.CompanyUUID,
		Sender:      emailData.Sender,
		Recipient:   emailData.Recipient,
		Subject:     emailData.Subject,
		Status:      emailData.Status,
		Source:      emailData.Source,
		Website:     emailData.Website,
		Payload:     emailData.Payload,
	}

	// Save the email to the database
	if err := initializers.DB.Create(&newEmail).Error; err != nil {
		// If an error occurs while saving the email, return an error response
		response.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	// Return the success response with the newly created email data
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": newEmail})
}


func GetEmails(c *gin.Context) {
	var emails []models.Email

	// Retrieve emails from the database
	if err := initializers.DB.Find(&emails).Error; err != nil {
		// If there's an error querying the database, return a 500 response
		response.WriteError(c, http.StatusInternalServerError, nil)
		return
	}

	// If no emails are found, return a 404 response
	if len(emails) == 0 {
		response.WriteError(c, http.StatusNotFound, nil)
		return
	}

	// Return emails in the response
	c.JSON(http.StatusOK, gin.H{"success": true, "data": emails})
}


func GetEmailByUUID(c *gin.Context) {
	// Declare a variable of type Email to hold the result
	var email models.Email
	
	// Get the email by UUID from the database
	if err := initializers.DB.Where("uuid = ?", c.Param("uuid")).First(&email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// If the email is not found, return a 404 response
			response.WriteError(c, http.StatusNotFound, nil)
		} else {
			// Handle any other error
			response.WriteError(c, http.StatusInternalServerError, nil)
		}
		return
	}

	// Return the email data in the response
	c.JSON(http.StatusOK, gin.H{"success": true, "data": email})
}


func UpdateEmailByUUID(c *gin.Context) {
	// Initialize validator
	var validate = validator.New()
	models.RegisterCustomValidations(validate) // Register custom validations

	// Bind and validate the request body
	var updateData models.Email
	if err := c.ShouldBindJSON(&updateData); err != nil {
		response.WriteError(c, http.StatusBadRequest, err)
		return
	}

	// Find the existing email record
	var email models.Email
	if err := initializers.DB.Where("uuid = ?", c.Param("uuid")).First(&email).Error; err != nil {
		response.WriteError(c, http.StatusNotFound, errors.New("email not found"))
		return
	}

	// Update only the provided fields
	if updateData.CompanyUUID != uuid.Nil {
		email.CompanyUUID = updateData.CompanyUUID
	}
	if updateData.Sender != "" {
		email.Sender = updateData.Sender
	}
	if updateData.Recipient != "" {
		email.Recipient = updateData.Recipient
	}
	if updateData.Subject != "" {
		email.Subject = updateData.Subject
	}
	if updateData.Status != "" {
		// Validate the status
		if !updateData.Status.IsValid() {
			response.WriteError(c, http.StatusBadRequest, errors.New("invalid status value"))
			return
		}
		email.Status = updateData.Status
	}
	if updateData.Source != "" {
		// Validate the source
		if !updateData.Source.IsValid() {
			response.WriteError(c, http.StatusBadRequest, errors.New("invalid source value"))
			return
		}
		email.Source = updateData.Source
	}
	if updateData.Website != "" {
		// Validate the website
		if !updateData.Website.IsValid() {
			response.WriteError(c, http.StatusBadRequest, errors.New("invalid website value"))
			return
		}
		email.Website = updateData.Website
	}
	if updateData.Payload != "" {
		email.Payload = updateData.Payload
	}

	// Save the updated email
	if err := initializers.DB.Save(&email).Error; err != nil {
		response.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	// Return the success response with the updated email data
	c.JSON(http.StatusOK, gin.H{"success": true, "data": email})
}


func DeleteEmailByUUID(c *gin.Context) {
	// Declare a variable of type Email to hold the result
	var email models.Email

	// Get the email by UUID from the database
	if err := initializers.DB.Where("uuid = ?", c.Param("uuid")).First(&email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// If the email is not found, return a 404 response
			response.WriteError(c, http.StatusNotFound, nil)
		} else {
			// Handle any other error
			response.WriteError(c, http.StatusInternalServerError, nil)
		}
		return
	}

	// Delete the email from the database
	if err := initializers.DB.Delete(&email).Error; err != nil {
		// If an error occurs while deleting the email, return an error response
		response.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, gin.H{"success": true})
}
