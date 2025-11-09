package auth

import (
	"fmt"
	"go-sse/user"
	"go-sse/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LoginHandler(c *gin.Context) {
	var loginDTO LoginDTO
	c.ShouldBindJSON(&loginDTO)
	validationErrors := util.Validate(loginDTO)
	if len(validationErrors) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Validation Error",
			"errors":  validationErrors,
		})
		return
	}

	authUser, err := user.GetUserByUsername(loginDTO.Username)
	if err != nil {
		util.ComparePassword("some password", loginDTO.Password)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "username or password is incorrect",
		})
		return
	}

	correctPassword, err := util.ComparePassword(*authUser.Password, loginDTO.Password)
	if err != nil {
		c.Error(err)
		return
	}
	if !correctPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "username or password is incorrect",
		})
		return
	}

	token, err := util.CreateJWT(authUser.Id.String())
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "cannot create token",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func RegisterHandler(c *gin.Context) {
	var registerDTO RegisterDTO
	c.ShouldBindJSON(&registerDTO)
	validationErrors := util.Validate(registerDTO)
	trimmedUsername := strings.TrimSpace(registerDTO.Username)
	fmt.Println(trimmedUsername)
	if usernameError := util.FindSlice(&validationErrors, func(validationError *util.ValidationError) bool {
		return validationError.Field == "username"
	}); usernameError == nil {
		_, err := user.GetUserByUsername(trimmedUsername)
		if err != nil && err != gorm.ErrRecordNotFound {
			c.Error(err)
			return
		} else if err == nil {
			validationErrors = append(validationErrors, util.ValidationError{
				Field:   "username",
				Tag:     "not_unique",
				Value:   trimmedUsername,
				Message: "username already exist",
			})
		}
	}
	if len(validationErrors) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request",
			"code":    "VALIDATION_ERROR",
			"errors":  validationErrors,
		})
		return
	}

	newUser, err := user.CreateUser(user.CreateUserDTO{
		Name:     registerDTO.Name,
		Username: &trimmedUsername,
		Password: &registerDTO.Password,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Register",
		"data":    newUser,
	})
}

func UserHandler(c *gin.Context) {
	authUser, ok := c.Get("authUser")

	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "auth user",
		"data":    authUser,
	})
}
