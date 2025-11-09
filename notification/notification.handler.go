package notification

import (
	"fmt"
	"go-sse/user"
	"go-sse/util"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetNotificationsHandler(c *gin.Context) {
	authUserContext, ok := c.Get("authUser")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}
	authUser, _ := authUserContext.(user.User)

	var getNotificationsDto GetNotificationsDto
	c.ShouldBind(&getNotificationsDto)
	validationErrors := util.Validate(getNotificationsDto)
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request",
			"code":    "VALIDATION_ERROR",
			"errors":  validationErrors,
		})
		return
	}

	data, count, err := GetPaginatedNotifications(authUser.Id, GetNotificationsDto{
		Page:  getNotificationsDto.Page,
		Limit: getNotificationsDto.Limit,
	})

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "get notifications",
		"metadata": gin.H{
			"total_data": count,
			"total_page": math.Ceil(float64(count) / float64(getNotificationsDto.Limit)),
		},
		"data": data,
	})
}

func StreamHandler(c *gin.Context) {
	authUserContext, ok := c.Get("authUser")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}
	authUser, _ := authUserContext.(user.User)

	newUuid, _ := uuid.NewRandom()
	client := &Client{
		Id:      newUuid,
		channel: make(chan string),
		done:    make(chan struct{}),
	}
	broker.AddClient(authUser.Id, client)

	w := c.Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		c.String(http.StatusInternalServerError, "streaming unsupported")
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	fmt.Fprintf(w, "event: connected\ndata: welcome %s\n\n", authUser.Name)
	flusher.Flush()

	for {
		select {
		case msg := <-client.channel:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()

		case <-c.Request.Context().Done():
			fmt.Printf("client %s disconnected", client.Id)
			close(client.done)
			return
		}
	}
}

func NotifyHandler(c *gin.Context) {
	userIdParam := c.Param("userId")
	userId, _ := uuid.Parse(userIdParam)
	message := c.PostForm("message")

	broker := GetBroker()
	broker.mu.RLock()
	clients, ok := broker.clients[userId]
	broker.mu.RUnlock()

	if ok {
		for _, client := range clients {
			client.channel <- message
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "notification sent",
	})
}
