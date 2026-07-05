package ws

import (
	"context"
	"log"
	"time"
)

type Room struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Members map[string]bool `json:"-"`
}

type Hub struct {
	ActiveClients map[string]*Client
	Rooms         map[string]*Room
	Register      chan *Client
	Unregister    chan *Client
	Broadcast     chan *Message
	Repo          Repository
}

func NewHub(repository Repository) *Hub {
	return &Hub{
		ActiveClients: make(map[string]*Client),
		Rooms:         make(map[string]*Room),
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
		Broadcast:     make(chan *Message, 5),
		Repo:          repository,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			// Register client globally by their UserID
			if _, ok := h.ActiveClients[cl.ID]; !ok {
				h.ActiveClients[cl.ID] = cl
				go h.Repo.SetUserOnlineStatus(context.Background(), cl.ID, true)
			}
		case cl := <-h.Unregister:
			// Remove client from global online users
			if _, ok := h.ActiveClients[cl.ID]; ok {
				delete(h.ActiveClients, cl.ID)
				close(cl.Message)
				go h.Repo.SetUserOnlineStatus(context.Background(), cl.ID, false)
			}
		case msg := <-h.Broadcast:
			// Look up room members in our cached map
			if room, ok := h.Rooms[msg.RoomID]; ok {
				// Broadcast to all online members of this room
				for userID := range room.Members {
					if cl, isOnline := h.ActiveClients[userID]; isOnline {
						cl.Message <- msg
					}
				}

				// Asynchronously write message to database
				// TODO: Implement a write pool on database to prevent bottleneck
				go func(msg *Message) {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					if err := h.Repo.WriteMessage(ctx, msg); err != nil {
						log.Printf("Error writing message to database: %v", err)
					}
				}(msg)
			}
		}
	}
}

// Shutdown safely disconnects all clients connected to this specific server node
// and ensures their database status is set back to offline to prevent zombies.
func (h *Hub) Shutdown() {
	for userID, client := range h.ActiveClients {
		h.Repo.SetUserOnlineStatus(context.Background(), userID, false)
		client.Conn.Close()
	}
}
