package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Enum interface with the IsValid method
type Enum interface {
    IsValid() bool
}

// Enum definitions for Status
type Status string

const (
	Sent   Status = "SENT"
	Failed Status = "FAILED"
)

func (s Status) IsValid() bool {
	switch s {
		case Sent, Failed:
		return true
	}
	return false
}

// Enum definitions for Source
type Source string

const (
	TrialCreated       		Source = "TRIAL_CREATED"
	TrialExpired       		Source = "TRIAL_EXPIRED"
	SubscriptionCreated 	Source = "SUBSCRIPTION_CREATED"
	SubscriptionRenewed 	Source = "SUBSCRIPTION_RENEWED"
	SubscriptionCancelled 	Source = "SUBSCRIPTION_CANCELLED"
	AccountCreation    		Source = "ACCOUNT_CREATION"
	ResetPassword      		Source = "RESET_PASSWORD"
	ChangeEmail        		Source = "CHANGE_EMAIL"
	DeleteAccount      		Source = "DELETE_ACCOUNT"
)

func (s Source) IsValid() bool {
	switch s {
		case TrialCreated, TrialExpired, SubscriptionCreated,
		SubscriptionRenewed, SubscriptionCancelled,
		AccountCreation, ResetPassword, ChangeEmail, DeleteAccount:
		return true
	}
	return false
}

// Enum definitions for Website
type Website string

const (
	IK  Website = "IK"
	MYE Website = "MYE"
	AK  Website = "AK"
)

func (w Website) IsValid() bool {
	switch w {
		case IK, MYE, AK:
		return true
	}
	return false
}

type Email struct {
	gorm.Model
	UUID        uuid.UUID `json:"uuid" gorm:"primaryKey;unique;"`
	CompanyUUID uuid.UUID `json:"company_uuid" gorm:"index" validate:"required"`
	Sender      string    `json:"sender" validate:"required,email"`
	Recipient   string    `json:"receiver" validate:"required,email"`
	Subject     string    `json:"subject" validate:"required"`
	Status      Status    `json:"status" validate:"required,status"`  // Custom validation for Status
	Source      Source    `json:"source" validate:"required,source"`  // Custom validation for Source
	Website     Website   `json:"website" validate:"required,website"` // Custom validation for Website
	Payload     string    `json:"payload" validate:"required"`
}

func (e *Email) BeforeCreate(tx *gorm.DB) (err error) {
	// Automatically set the UUID before creating the record
	if e.UUID == uuid.Nil {
		e.UUID = uuid.New()
	}
	return nil
}

// Register custom validations
func RegisterCustomValidations(v *validator.Validate) {
	v.RegisterValidation("status", ValidateStatus)
	v.RegisterValidation("source", ValidateSource)
	v.RegisterValidation("website", ValidateWebsite)
}

// ValidateStatus validates the status field to ensure it's one of the valid enum values
func ValidateStatus(fl validator.FieldLevel) bool {
	status, ok := fl.Field().Interface().(Status)
	if !ok {
		return false
	}
	return status.IsValid()
}

// ValidateSource validates the source field to ensure it's one of the valid enum values
func ValidateSource(fl validator.FieldLevel) bool {
	source, ok := fl.Field().Interface().(Source)
	if !ok {
		return false
	}
	return source.IsValid()
}

// ValidateWebsite validates the website field to ensure it's one of the valid enum values
func ValidateWebsite(fl validator.FieldLevel) bool {
	website, ok := fl.Field().Interface().(Website)
	if !ok {
		return false
	}
	return website.IsValid()
}

