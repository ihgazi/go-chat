package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ihgazi/go-chat/internal/user"
	"github.com/ihgazi/go-chat/internal/ws"
	"github.com/ihgazi/go-chat/util"
)

// GIN Router to create API endpoints

var r *gin.Engine

func Init(userHandler *user.Handler, wsHandler *ws.Handler) *gin.Engine {
	r = gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	r.POST("/signup", userHandler.CreateUser)
	r.POST("/login", userHandler.Login)
	r.GET("/logout", userHandler.Logout)

	// Protected routes
	protected := r.Group("/ws")
	protected.Use(util.JWTValidateToken())
	{
		protected.GET("/auth", userHandler.AuthUser)
		protected.POST("/createRoom", wsHandler.CreateRoom)
		protected.GET("/getRooms", wsHandler.GetRooms)
		protected.GET("/getClients/:roomId", wsHandler.GetClients)
	}

    // Moved joinRoom out of protected group
    // for WebSocket issues in deployment
	r.GET("/joinRoom/:roomId", wsHandler.JoinRoom)

	return r
}
