package main

import (
	"fmt"

	"github.com/farhan-nahid/email-service/initializers"
	"github.com/farhan-nahid/email-service/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
}

func main() {
	fmt.Println("Attempting to Migration")
	err:= initializers.DB.AutoMigrate(&models.Email{})

	if err != nil {
		fmt.Println("Migration Failed")
		panic(err)
	}

	fmt.Println("Migration Successful")
}