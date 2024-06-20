package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ihgazi/go-chat/internal/user"
	"github.com/ihgazi/go-chat/internal/ws"
)

// GIN Router to create API endpoints

var r *gin.Engine

func Init(userHandler *user.Handler, wsHandler *ws.Handler) *gin.Engine {
	r = gin.Default()

	r.POST("/signup", userHandler.CreateUser)
	r.POST("/login", userHandler.Login)
	r.GET("/logout", userHandler.Logout)

	r.POST("/ws/createRoom", wsHandler.CreateRoom)
	r.GET("/ws/joinRoom/:roomId", wsHandler.JoinRoom)
	r.GET("/ws/getRooms", wsHandler.GetRooms)
	r.GET("/ws/getClients/:roomId", wsHandler.GetClients)

	return r
}
