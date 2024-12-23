package main

import (
	"github.com/farhan-nahid/email-service/initializers"
	"github.com/farhan-nahid/email-service/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
}

func main() {
	initializers.DB.AutoMigrate(&models.Email{})
}