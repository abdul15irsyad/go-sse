package auth

import "github.com/gin-gonic/gin"

func AuthRoute(router *gin.Engine) {
	auth := router.Group("/auth")
	auth.POST("/login", LoginHandler)
	auth.POST("/register", RegisterHandler)
	auth.Use(AuthMiddleware)
	auth.POST("/logout", LogoutHandler)
	auth.GET("/user", UserHandler)
}
