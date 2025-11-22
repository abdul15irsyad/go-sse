package notification

import (
	"encoding/json"
	"fmt"
	"go-sse/user"
	"net/http"
	"strconv"

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
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	data, count, err := GetPaginatedNotifications(authUser.Id, page, limit, status)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "get notifications",
		"count":   count,
		"data":    data,
	})
}

func GetCountNotificationsHandler(c *gin.Context) {
	authUserContext, _ := c.Get("authUser")
	authUser, _ := authUserContext.(user.User)
	status := c.Query("status")

	count, err := GetCountNotifications(authUser.Id, status)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "get " + status + "notifications",
		"count":   count,
	})
}

func ReadNotificationHandler(c *gin.Context) {
	authUserContext, _ := c.Get("authUser")
	authUser, _ := authUserContext.(user.User)
	idParam := c.Param("id")
	id, _ := uuid.Parse(idParam)

	err := ReadNotification(id, authUser.Id)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "notification read",
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
		channel: make(chan Notification),
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
		case notification := <-client.channel:
			fmt.Println("new notification", notification)
			notificationJson, _ := json.Marshal(notification)
			fmt.Println("notification", string(notificationJson))
			fmt.Fprintf(w, "data: %s\n\n", string(notificationJson))
			flusher.Flush()

		case <-c.Request.Context().Done():
			fmt.Printf("client %s disconnected", client.Id)
			close(client.done)
			return
		}
	}
}

func PokeHandler(c *gin.Context) {
	authUserContext, ok := c.Get("authUser")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}
	authUser, _ := authUserContext.(user.User)

	frendIdParam := c.Param("frendId")
	frendId, _ := uuid.Parse(frendIdParam)

	broker := GetBroker()
	broker.mu.RLock()
	clients, ok := broker.clients[frendId]
	broker.mu.RUnlock()

	notification, err := CreateNotification(CreateNotificationDTO{
		UserId:  frendId,
		Title:   "Friend Poke",
		Message: fmt.Sprintf("you just poked by %s", authUser.Name),
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	if ok {
		for _, client := range clients {
			fmt.Println("client", client)
			fmt.Println("*client", *client)
			fmt.Println("notification", notification)
			client.channel <- notification
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "notification sent",
	})
}
