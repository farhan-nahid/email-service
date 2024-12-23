package routes

import (
	"github.com/farhan-nahid/email-service/controllers"
	"github.com/farhan-nahid/email-service/middleware"
	"github.com/farhan-nahid/email-service/models"
	"github.com/gin-gonic/gin"
)

func EmailRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.POST("/email", middleware.BindAndValidate[models.Email](), controllers.CreateEmail)
		v1.GET("/email", controllers.GetEmails)
		v1.GET("/email/deleted", controllers.GetDeletedEmails)
		v1.GET("/email/:uuid", controllers.GetEmailByUUID)
		v1.PATCH("/email/:uuid", controllers.UpdateEmailByUUID)
		v1.DELETE("/email/:uuid", controllers.DeleteEmailByUUID)
	}
}