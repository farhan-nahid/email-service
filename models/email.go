package models

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Enum definitions for Status
type Status string

const (
	Sent     Status = "SENT"
	Failed   Status = "FAILED"
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
	TrialCreate    	    Source = "TRIAL_CREATED"
	TrialExpire    	    Source = "TRIAL_EXPIRED"
	SubscriptionCreate  Source = "SUBSCRIPTION_CREATED"
	SubscriptionRenew   Source = "SUBSCRIPTION_RENEWED"
	SubscriptionCancel  Source = "SUBSCRIPTION_CANCELLED"
	AccountCreation 	Source = "ACCOUNT_CREATION"
	ResetPassword   	Source = "RESET_PASSWORD"
	ChangeEmail     	Source = "CHANGE_EMAIL"
	DeleteAccount   	Source = "DELETE_ACCOUNT"
)

func (s Source) IsValid() bool {
	switch s {
		case TrialCreate, TrialExpire, SubscriptionCreate, SubscriptionCancel, SubscriptionRenew, AccountCreation, ResetPassword, ChangeEmail, DeleteAccount:
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
	UUID     	uuid.UUID `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();primaryKey;unique;"`
	CompanyUUID uuid.UUID `json:"company_uuid" gorm:"index"`
	Sender   	string    `json:"sender"`
	Recipient 	string    `json:"receiver"`
	Subject  	string    `json:"subject"`
	Status   	Status    `json:"status"`   // SENT, FAILED
	Source   	Source    `json:"source"`   // TRIAL_CREATED, TRIAL_EXPIRED, SUBSCRIPTION_CREATED, SUBSCRIPTION_RENEWED, SUBSCRIPTION_CANCELLED, ACCOUNT_CREATION, RESET_PASSWORD, CHANGE_EMAIL, DELETE_ACCOUNT
	Website  	Website   `json:"website"` 	// IK, MYE, AK
	Payload  	string    `json:"payload"`
}

// BeforeSave hook to validate enums
func (e *Email) BeforeSave(tx *gorm.DB) (err error) {
	if !e.Source.IsValid() {
		return errors.New("invalid source value")
	}
	if !e.Website.IsValid() {
		return errors.New("invalid website value")
	}
	return nil
}