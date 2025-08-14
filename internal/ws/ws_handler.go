package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		s,
	}
}

func (h *Handler) CreateRoom(c *gin.Context) {
	var req CreateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Service.CreateRoom(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, res)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// TODO: whitelist the frontend origin
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) JoinRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	roomID := c.Param("roomId")
	clientIDStr := c.Query("userID")
	usernameStr := c.Query("username")

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientIDStr,
		RoomID:   roomID,
		Username: usernameStr,
	}

	m := &Message{
		Content:  "A new user has joined the room",
		RoomID:   roomID,
		Username: usernameStr,
		UserID:   clientIDStr,
	}

	err = h.Service.JoinRoom(c.Request.Context(), cl, m)
	if err != nil {
		log.Fatalf(err.Error())
		conn.Close()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join room!"})
	}
}

// get currently active rooms in hub
func (h *Handler) GetRooms(c *gin.Context) {
	r := h.Service.GetRooms(c.Request.Context())

	c.JSON(http.StatusOK, r)
}

func (h *Handler) GetClients(c *gin.Context) {
	roomId := c.Param("roomId")

	clients := h.Service.GetClients(c.Request.Context(), roomId)

	c.JSON(http.StatusOK, clients)
}
