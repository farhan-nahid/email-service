package models

import (
	"errors"
	"net/mail"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Enum interface with the IsValid method
type Enum interface {
	IsValid() bool
}

// Custom EmailAddress type
type EmailAddress string

// IsValid checks if the EmailAddress is valid
func (e EmailAddress) IsValid() bool {
	_, err := mail.ParseAddress(string(e))
	return err == nil
}

// Scan implements the sql.Scanner interface for database compatibility
func (e *EmailAddress) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("failed to scan Email Address: value is not a string")
	}
	*e = EmailAddress(str)
	return nil
}

// Value implements the driver.Valuer interface for database compatibility
func (e EmailAddress) Value() (interface{}, error) {
	return string(e), nil
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
	TrialCreated          Source = "TRIAL_CREATED"
	TrialExpired          Source = "TRIAL_EXPIRED"
	SubscriptionCreated   Source = "SUBSCRIPTION_CREATED"
	SubscriptionRenewed   Source = "SUBSCRIPTION_RENEWED"
	SubscriptionCancelled Source = "SUBSCRIPTION_CANCELLED"
	AccountCreation       Source = "ACCOUNT_CREATION"
	ResetPassword         Source = "RESET_PASSWORD"
	ChangeEmail           Source = "CHANGE_EMAIL"
	DeleteAccount         Source = "DELETE_ACCOUNT"
)

func (s Source) IsValid() bool {
	switch s {
		case TrialCreated, TrialExpired, SubscriptionCreated, SubscriptionRenewed, 
		SubscriptionCancelled, AccountCreation, ResetPassword, ChangeEmail, DeleteAccount:
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

// Email model
type Email struct {
	gorm.Model
	UUID        uuid.UUID    `json:"uuid" gorm:"primaryKey;unique;"`
	CompanyUUID uuid.UUID    `json:"company_uuid" gorm:"index" validate:"required,uuid"`
	Sender      EmailAddress `json:"sender" validate:"required,email_address"`
	Recipient   EmailAddress `json:"receiver" validate:"required,email_address"`
	Subject     string       `json:"subject" validate:"required"`
	Status      Status       `json:"status" validate:"required,status"`
	Source      Source       `json:"source" validate:"required,source"`
	Website     Website      `json:"website" validate:"required,website"`
	Payload     string       `json:"payload" validate:"required"`
}

// BeforeCreate hook to set UUID automatically
func (e *Email) BeforeCreate(tx *gorm.DB) (err error) {
	if e.UUID == uuid.Nil {
		e.UUID = uuid.New()
	}
	return nil
}

// IsValid validates the Email struct fields
func (e *Email) IsValid() bool {
	return e.Sender.IsValid() &&
		e.Recipient.IsValid() &&
		e.Status.IsValid() &&
		e.Source.IsValid() &&
		e.Website.IsValid() &&
		e.UUIDIsValid()
}

// UUIDIsValid validates UUID fields
func (e *Email) UUIDIsValid() bool {
	return e.UUID != uuid.Nil && e.CompanyUUID != uuid.Nil
}

// Register custom validations
func RegisterCustomValidations(v *validator.Validate) {
	v.RegisterValidation("status", ValidateStatus)
	v.RegisterValidation("source", ValidateSource)
	v.RegisterValidation("website", ValidateWebsite)
	v.RegisterValidation("uuid", ValidateUUID)
	v.RegisterValidation("email_address", ValidateEmailAddress)
}

// ValidateStatus validates the status field
func ValidateStatus(fl validator.FieldLevel) bool {
	status, ok := fl.Field().Interface().(Status)
	return ok && status.IsValid()
}

// ValidateSource validates the source field
func ValidateSource(fl validator.FieldLevel) bool {
	source, ok := fl.Field().Interface().(Source)
	return ok && source.IsValid()
}

// ValidateWebsite validates the website field
func ValidateWebsite(fl validator.FieldLevel) bool {
	website, ok := fl.Field().Interface().(Website)
	return ok && website.IsValid()
}

// ValidateUUID validates a UUID
func ValidateUUID(fl validator.FieldLevel) bool {
	id, ok := fl.Field().Interface().(uuid.UUID)
	return ok && id != uuid.Nil
}

// ValidateEmailAddress validates the EmailAddress type
func ValidateEmailAddress(fl validator.FieldLevel) bool {
	email, ok := fl.Field().Interface().(EmailAddress)
	return ok && email.IsValid()
}
