package ws

import (
	"fmt"
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

// Connect upgrades the HTTP connection to a WebSocket session
func (h *Handler) Connect(c *gin.Context) {
	clientIDStr := c.GetString("userID")
	usernameStr := c.GetString("username")

	// Validate query parameters before upgrading the connection
	if clientIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID query parameter is required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message),
		ID:       clientIDStr,
		Username: usernameStr,
	}

	err = h.Service.Connect(c.Request.Context(), cl)
	if err != nil {
		log.Printf("error connecting client: %v", err)
		conn.Close()
	}
}

// JoinRoom is a REST endpoint to add a user to a room
func (h *Handler) JoinRoom(c *gin.Context) {
	roomID := c.Param("roomId")
	clientIDStr := c.GetString("userID")

	// Add validation to prevent database casting errors
	if clientIDStr == "" || roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID query parameter is required"})
		return
	}

	err := h.Service.JoinRoom(c.Request.Context(), roomID, clientIDStr)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to join room %s: %v", roomID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully joined room"})
}

// Get list of users in the Room
func (h *Handler) GetRooms(c *gin.Context) {
	r := h.Service.GetRooms(c.Request.Context())

	c.JSON(http.StatusOK, r)
}

// GetMyRooms returns only the rooms the current user is a member of
func (h *Handler) GetMyRooms(c *gin.Context) {
	userIDStr := c.GetString("userID")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID query parameter is required"})
		return
	}

	rooms := h.Service.GetMyRooms(c.Request.Context(), userIDStr)
	c.JSON(http.StatusOK, rooms)
}

func (h *Handler) GetClients(c *gin.Context) {
	roomId := c.Param("roomId")

	clients := h.Service.GetClients(c.Request.Context(), roomId)

	c.JSON(http.StatusOK, clients)
}
