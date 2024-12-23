package controllers

import (
	"github.com/farhan-nahid/email-service/initializers"
	"github.com/farhan-nahid/email-service/models"
	"github.com/gin-gonic/gin"
)

func CreateEmail(c *gin.Context) {
	// parse request body
	email := models.Email{
		Subject: c.PostForm("subject"),
		Sender: c.PostForm("sender"),
		Recipient: c.PostForm("recipient"),
		Payload: c.PostForm("body"),
		Website: models.Website(c.PostForm("website")),
		Source: models.Source(c.PostForm("source")),
	}

	result := initializers.DB.Create(&email) // pass pointer of data to Create


	
	// validate request body
	// create email
	// send email
	// return response
}