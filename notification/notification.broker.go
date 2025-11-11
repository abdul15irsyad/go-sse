package notification

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type Client struct {
	Id      uuid.UUID
	channel chan Notification
	done    chan struct{}
}

type Broker struct {
	mu      sync.RWMutex
	clients map[uuid.UUID][]*Client
}

var broker = Broker{
	clients: make(map[uuid.UUID][]*Client),
}

func GetBroker() *Broker {
	return &broker
}

func (b *Broker) SendNotificationToUser(userId uuid.UUID, notification Notification) {
	b.mu.RLock()
	clients, ok := b.clients[userId]
	b.mu.RUnlock()

	if !ok || len(clients) == 0 {
		fmt.Printf("no active connection for user %s\n", userId)
		return
	}

	for _, c := range clients {
		select {
		case c.channel <- notification:
		default:
			fmt.Printf("client %s channel full, dropping message\n", c.Id)
		}
	}
}

func (b *Broker) AddClient(userId uuid.UUID, client *Client) {
	b.mu.Lock()
	b.clients[userId] = append(b.clients[userId], client)
	b.mu.Unlock()

	fmt.Printf("client %s (user %s) added\n", client.Id, userId)

	go func() {
		<-client.done
		b.RemoveClient(userId, client.Id)
	}()
}

func (b *Broker) RemoveClient(userId uuid.UUID, clientId uuid.UUID) {
	b.mu.Lock()
	defer b.mu.Unlock()

	clients := b.clients[userId]
	filtered := make([]*Client, 0, len(clients))
	for _, c := range clients {
		if c.Id != clientId {
			filtered = append(filtered, c)
		}
	}

	if len(filtered) > 0 {
		b.clients[userId] = filtered
	} else {
		delete(b.clients, userId)
	}

	fmt.Printf("client %s (user %s) removed\n", clientId, userId)
}
