package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/farhan-nahid/email-service/initializers"
	"github.com/farhan-nahid/email-service/models"
	"github.com/farhan-nahid/email-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// checks the given string is a valid UUID
func isValidUUID(str string) bool {
	// Attempt to parse the string as a UUID
	_, err := uuid.Parse(str)
	return err == nil
}

func CreateEmail(c *gin.Context) {
	validatedData, exists := c.Get("validatedData")
	if !exists {
		utils.ErrorResponse(c, http.StatusBadRequest, errors.New("data not found in context"))
		return
	}

	emailData, ok := validatedData.(models.Email)
	if !ok {
		utils.ErrorResponse(c, http.StatusBadRequest, errors.New("invalid data format"))
		return
	}

	var payload interface{}

	// Parse (unmarshal) the JSON string into the struct
	err := json.Unmarshal([]byte(emailData.Payload), &payload)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, errors.New("invalid payload format"))
	}

	// Create a new email instance using the validated data
	newEmail := models.Email{
		Name: 	 	 emailData.Name,
		CompanyUUID: emailData.CompanyUUID,
		Sender:      emailData.Sender,
		Recipient:   emailData.Recipient,
		Subject:     emailData.Subject,
		Source:      emailData.Source,
		Website:     emailData.Website,
		Payload:     emailData.Payload,
		Status:      "SENT",
	}

	var sender string;

	if emailData.Website == "AK" {
		sender = "Attendance Keeper <" + string(emailData.Sender) + ">"
	} else if emailData.Website == "MYE" {
		sender = "Manage Your Ecommerce <" + string(emailData.Sender) + ">"
	} else {
		sender = "Inventory Keeper <" + string(emailData.Sender) + ">"
	}


	// Send Email
	err = utils.SendEmail(utils.Data{
		Name: emailData.Name,
		Sender: sender,
		Receiver: string(emailData.Recipient),
		Subject: emailData.Subject,
		Payload: payload,
	}, "/templates/" + string(emailData.Website) + "/" + string(emailData.Source) + ".html")

	
	if err !=  nil{
		newEmail.Status = "FAILED"
		// Save the email to the database
		if err := initializers.DB.Create(&newEmail).Error; err != nil {
			// If an error occurs while saving the email, return an error utils
			utils.ErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	// Save the email to the database
	if err := initializers.DB.Create(&newEmail).Error; err != nil {
		// If an error occurs while saving the email, return an error utils
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	// Return a success utils
	utils.SuccessResponse(c, http.StatusCreated, newEmail, "Email created successfully")
}


func GetEmails(c *gin.Context) {
	var emails []models.Email

	// Retrieve emails from the database
	if err := initializers.DB.Find(&emails).Error; err != nil {
		// If there's an error querying the database, return a 500 utils
		utils.ErrorResponse(c, http.StatusInternalServerError, nil)
		return
	}

	// If no emails are found, return a 404 utils
	if len(emails) == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, nil)
		return
	}

	// Return a success utils
	utils.SuccessResponse(c, http.StatusOK, emails, "Emails retrieved successfully")
}


func GetDeletedEmails(c *gin.Context) {
	var emails []models.Email

	// Retrieve deleted emails from the database
	if err := initializers.DB.Unscoped().Where("deleted_at IS NOT NULL").Find(&emails).Error; err != nil {
		// If there's an error querying the database, return a 500 utils
		utils.ErrorResponse(c, http.StatusInternalServerError, nil)
		return
	}

	// If no emails are found, return a 404 utils
	if len(emails) == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, nil)
		return
	}

	// Return a success utils
	utils.SuccessResponse(c, http.StatusOK, emails, "Deleted emails retrieved successfully")
}


func GetEmailByUUID(c *gin.Context) {
	// Validate the UUID	
	if !isValidUUID(c.Param("uuid")) {
		utils.ErrorResponse(c, http.StatusBadRequest, errors.New("invalid UUID format in request URL"))
		return
	}

	// Declare a variable of type Email to hold the result
	var email models.Email

	
	
	// Get the email by UUID from the database
	if err := initializers.DB.Where("uuid = ?", c.Param("uuid")).First(&email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// If the email is not found, return a 404 utils
			utils.ErrorResponse(c, http.StatusNotFound, nil)
		} else {
			// Handle any other error
			utils.ErrorResponse(c, http.StatusInternalServerError, nil)
		}
		return
	}

	// Return a success utils
	utils.SuccessResponse(c, http.StatusOK, email, "Email retrieved successfully")
}


func UpdateEmailByUUID(c *gin.Context) {
	// Validate the UUID	
	if !isValidUUID(c.Param("uuid")) {
		utils.ErrorResponse(c, http.StatusBadRequest, errors.New("invalid UUID format in request URL"))
		return
	}

	// Initialize validator
	var validate = validator.New()
	models.RegisterCustomValidations(validate) // Register custom validations

	// Bind and validate the request body
	var updateData models.Email
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	// Find the existing email record
	var email models.Email
	if err := initializers.DB.Where("uuid = ?", c.Param("uuid")).First(&email).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, errors.New("email not found"))
		return
	}

	if updateData.Name != "" {
		email.Name = updateData.Name
	}

	// Update only the provided fields
	if updateData.CompanyUUID != uuid.Nil {
		email.CompanyUUID = updateData.CompanyUUID
	}

	if updateData.Sender != "" {
		// Validate the Sender email address
		if !updateData.Sender.IsValid() {
			utils.ErrorResponse(c, http.StatusBadRequest, errors.New("invalid sender email address"))
			return
		}
		email.Sender = updateData.Sender
	}

	if updateData.Recipient != "" {
		// Validate the Recipient email address
		if !updateData.Recipient.IsValid() {
			utils.ErrorResponse(c, http.StatusBadRequest, errors.New("invalid recipient email address"))
			return
		}
		email.Recipient = updateData.Recipient
	}

	if updateData.Subject != "" {
		email.Subject = updateData.Subject
	}

	if updateData.Status != "" {
		// Validate the status
		if !updateData.Status.IsValid() {
			utils.ErrorResponse(c, http.StatusBadRequest, errors.New("invalid status value"))
			return
		}
		email.Status = updateData.Status
	}

	if updateData.Source != "" {
		// Validate the source
		if !updateData.Source.IsValid() {
			utils.ErrorResponse(c, http.StatusBadRequest, errors.New("invalid source value"))
			return
		}
		email.Source = updateData.Source
	}

	if updateData.Website != "" {
		// Validate the website
		if !updateData.Website.IsValid() {
			utils.ErrorResponse(c, http.StatusBadRequest, errors.New("invalid website value"))
			return
		}
		email.Website = updateData.Website
	}

	if updateData.Payload != "" {
		email.Payload = updateData.Payload
	}

	// Save the updated email
	if err := initializers.DB.Save(&email).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	// Return a success utils
	utils.SuccessResponse(c, http.StatusOK, email, "Email updated successfully")
}


func DeleteEmailByUUID(c *gin.Context) {
	// Validate the UUID format
	if !isValidUUID(c.Param("uuid")) {
		utils.ErrorResponse(c, http.StatusBadRequest, errors.New("invalid UUID format in request URL"))
		return
	}

	// Declare a variable of type Email to hold the result
	var email models.Email

	// Attempt to find the email by UUID in the database
	if err := initializers.DB.Where("uuid = ?", c.Param("uuid")).First(&email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// If the email is not found, return a 404 utils
			utils.ErrorResponse(c, http.StatusNotFound, errors.New("email not found"))
		} else {
			// Handle other database errors
			utils.ErrorResponse(c, http.StatusInternalServerError, err)
		}
		return
	}

	// If it is already deleted, return an error utils
	if email.DeletedAt.Valid {
		utils.ErrorResponse(c, http.StatusBadRequest, errors.New("email already deleted"))
		return
	}

	// Delete the email from the database
	if err := initializers.DB.Delete(&email).Error; err != nil {
		// If an error occurs while deleting, return an error utils
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	// Return a success utils after deletion
	utils.SuccessResponse(c, http.StatusOK, nil, "Email deleted successfully")
}
