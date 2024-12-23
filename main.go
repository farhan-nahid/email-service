package main

import (
	"github.com/farhan-nahid/email-service/controllers"
	"github.com/farhan-nahid/email-service/initializers"
	"github.com/gin-gonic/gin"
)

// init() is called automatically, no need to call it explicitly
func init() {
	initializers.LoadEnvVariables()
	// initializers.ConnectToDatabase()
}

func main() {
	r := gin.Default()

	r.GET("/health-check", controllers.HealthCheck)
	r.Run()
}