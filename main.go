package main

import (
	"fmt"
	"go-sse/auth"
	"go-sse/config"
	"go-sse/notification"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	auth.AuthRoute(r)
	notification.NotificationRoute(r)

	fmt.Println("server running on port:", config.Port)
	err := r.Run(fmt.Sprintf(":%d", config.Port))
	if err != nil {
		panic(err)
	}
}
