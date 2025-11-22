package notification

import (
	"go-sse/auth"

	"github.com/gin-gonic/gin"
)

func NotificationRoute(router *gin.Engine) {
	notification := router.Group("/notif")
	notification.Use(auth.AuthMiddleware)
	notification.GET("/", GetNotificationsHandler)
	notification.GET("/count", GetCountNotificationsHandler)
	notification.GET("/read/:id", ReadNotificationHandler)
	notification.GET("/stream", StreamHandler)
	notification.POST("/poke/:frendId", PokeHandler)
}
