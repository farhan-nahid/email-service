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

// ------------------- EmailAddress Type ------------------- //

type EmailAddress string

func (e EmailAddress) IsValid() bool {
	_, err := mail.ParseAddress(string(e))
	return err == nil
}

func (e *EmailAddress) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("failed to scan Email Address: value is not a string")
	}
	*e = EmailAddress(str)
	return nil
}

func (e EmailAddress) Value() (interface{}, error) {
	return string(e), nil
}

// ------------------- Enums ------------------- //

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

// ------------------- Email Model ------------------- //

type Email struct {
	gorm.Model
	UUID        uuid.UUID    `json:"uuid" gorm:"primaryKey;unique;"`
	CompanyUUID uuid.UUID    `json:"company_uuid" gorm:"index" validate:"required,uuid"`
	Name        string       `json:"name" validate:"required"`
	Sender      EmailAddress `json:"sender" validate:"required,email_address"`
	Recipient   EmailAddress `json:"receiver" validate:"required,email_address"`
	Subject     string       `json:"subject" validate:"required"`
	Status      Status       `json:"status" validate:"status"`
	Source      Source       `json:"source" validate:"required,source"`
	Website     Website      `json:"website" validate:"required,website"`
	Payload     string       `json:"payload" validate:"required,json"`
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
		e.Website.IsValid() 
}


// ------------------- Custom Validations ------------------- //

func RegisterCustomValidations(v *validator.Validate) {
	v.RegisterValidation("status", ValidateStatus)
	v.RegisterValidation("source", ValidateSource)
	v.RegisterValidation("website", ValidateWebsite)
	v.RegisterValidation("uuid", ValidateUUID)
	v.RegisterValidation("email_address", ValidateEmailAddress)
}

func ValidateStatus(fl validator.FieldLevel) bool {
	status, ok := fl.Field().Interface().(Status)
	return ok && status.IsValid()
}

func ValidateSource(fl validator.FieldLevel) bool {
	source, ok := fl.Field().Interface().(Source)
	return ok && source.IsValid()
}

func ValidateWebsite(fl validator.FieldLevel) bool {
	website, ok := fl.Field().Interface().(Website)
	return ok && website.IsValid()
}

func ValidateUUID(fl validator.FieldLevel) bool {
	id, ok := fl.Field().Interface().(uuid.UUID)
	return ok && id != uuid.Nil
}

func ValidateEmailAddress(fl validator.FieldLevel) bool {
	email, ok := fl.Field().Interface().(EmailAddress)
	return ok && email.IsValid()
}






// package models

// import (
// 	"github.com/google/uuid"
// 	"gorm.io/gorm"
// )

// // ------------------- Enums ------------------- //

// type Status string

// const (
// 	Sent   Status = "SENT"
// 	Failed Status = "FAILED"
// )

// type Source string

// const (
// 	TrialCreated          Source = "TRIAL_CREATED"
// 	TrialExpired          Source = "TRIAL_EXPIRED"
// 	SubscriptionCreated   Source = "SUBSCRIPTION_CREATED"
// 	SubscriptionRenewed   Source = "SUBSCRIPTION_RENEWED"
// 	SubscriptionCancelled Source = "SUBSCRIPTION_CANCELLED"
// 	AccountCreation       Source = "ACCOUNT_CREATION"
// 	ResetPassword         Source = "RESET_PASSWORD"
// 	ChangeEmail           Source = "CHANGE_EMAIL"
// 	DeleteAccount         Source = "DELETE_ACCOUNT"
// )

// type Website string

// const (
// 	IK  Website = "IK"
// 	MYE Website = "MYE"
// 	AK  Website = "AK"
// )

// // ------------------- Email Model ------------------- //

// type Email struct {
// 	gorm.Model
// 	UUID        uuid.UUID    `json:"uuid" gorm:"primaryKey;unique;not null"`
// 	CompanyUUID uuid.UUID    `json:"company_uuid" gorm:"index;not null" validate:"required,uuid"`
// 	Name        string       `json:"name" validate:"required"`
// 	Sender      string       `json:"sender" validate:"required,email"`
// 	Recipient   string       `json:"receiver" validate:"required,email"`
// 	Subject     string       `json:"subject" validate:"required"`
// 	Status      Status       `json:"status" gorm:"type:enum('SENT', 'FAILED');not null" validate:"required,enum_status"`
// 	Source      Source       `json:"source" gorm:"type:enum('TRIAL_CREATED', 'TRIAL_EXPIRED', 'SUBSCRIPTION_CREATED', 'SUBSCRIPTION_RENEWED', 'SUBSCRIPTION_CANCELLED', 'ACCOUNT_CREATION', 'RESET_PASSWORD', 'CHANGE_EMAIL', 'DELETE_ACCOUNT');not null" validate:"required,enum_source"`
// 	Website     Website      `json:"website" gorm:"type:enum('IK', 'AK', 'MYE');not null" validate:"required,enum_website"`
// 	Payload     string       `json:"payload" validate:"required,json"`
// }

// // BeforeCreate hook to set UUID automatically
// func (e *Email) BeforeCreate(tx *gorm.DB) (err error) {
// 	if e.UUID == uuid.Nil {
// 		e.UUID = uuid.New()
// 	}
// 	return nil
// }
