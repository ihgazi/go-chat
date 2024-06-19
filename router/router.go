package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ihgazi/go-chat/internal/user"
)

// GIN Router to create API endpoints

var r *gin.Engine

func Init(userHandler *user.Handler) *gin.Engine {
	r = gin.Default()

	r.POST("/signup", userHandler.CreateUser)
	r.POST("/login", userHandler.Login)
	r.GET("/logout", userHandler.Logout)

	return r
}
