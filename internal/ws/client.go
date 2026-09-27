package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       string `json:"id"`
	Username string `json:"username"`
}

type Message struct {
	Type     string `json:"type"`
	Event    string `json:"event,omitempty"`
	Content  string `json:"content"`
	RoomID   string `json:"room_id"`
	Username string `json:"username"`
	UserID   string `json:"user_id"`
}

// Take message from client channel for passing to frontend
func (cl *Client) WriteMessage() {
	defer func() {
		cl.Conn.Close()
	}()

	for {
		message, ok := <-cl.Message
		if !ok {
			return
		}

		cl.Conn.WriteJSON(message)
	}
}

// Read message from frontend
func (cl *Client) ReadMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- cl
		cl.Conn.Close()
	}()

	for {
		var msg Message
		err := cl.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Enforce the sender's identity to prevent spoofing
		msg.Username = cl.Username
		msg.UserID = cl.ID
		if msg.Type == "" {
			msg.Type = "user"
		}

		hub.Broadcast <- &msg
	}
}
