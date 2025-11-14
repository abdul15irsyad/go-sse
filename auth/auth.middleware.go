package auth

import (
	"errors"
	"fmt"
	"go-sse/user"
	"go-sse/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AuthMiddleware(c *gin.Context) {
	accessToken, err := c.Cookie(ACCESS_TOKEN_KEY)
	fmt.Println("accessToken", accessToken)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid credential",
		})
		return
	}

	sub, err := util.ParseJWT(accessToken)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "token expired",
				"code":    "TOKEN_EXPIRED",
			})
			return
		}
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid credential",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid credential",
		})
		return
	}

	// check to database
	userId, err := uuid.Parse(sub)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid credential",
		})
		return
	}
	authUser, err := user.GetUser(userId)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid credential",
			})
			return
		}
		c.Error(err)
		return
	}

	c.Set("authUser", authUser)
	c.Next()
}
