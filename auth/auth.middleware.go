package auth

import (
	"errors"
	"go-sse/user"
	"go-sse/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AuthMiddleware(c *gin.Context) {
	authorization := c.GetHeader("Authorization")

	if authorization == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid credential",
		})
		return
	}
	accessToken := strings.Split(authorization, " ")[1]

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
