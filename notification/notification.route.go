package notification

import (
	"go-sse/auth"

	"github.com/gin-gonic/gin"
)

func NotificationRoute(router *gin.Engine) {
	notification := router.Group("/notifications")
	notification.Use(auth.AuthMiddleware)
	notification.GET("/", GetNotificationsHandler)
	notification.GET("/stream", StreamHandler)
	notification.POST("/send/:userId", SendHandler)
}
