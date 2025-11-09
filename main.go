package main

import (
	"fmt"
	"go-sse/auth"
	"go-sse/config"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	auth.AuthRoute(r)

	fmt.Println("server running on port:", config.Port)
	err := r.Run(fmt.Sprintf(":%d", config.Port))
	if err != nil {
		panic(err)
	}
}
