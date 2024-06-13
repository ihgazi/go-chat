package router 

import (
    "github.com/ihgazi/go-chat/internal/user"
    "github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(userHandler *user.Handler) {
    r = gin.Default()

    r.POST("/signup", userHandler.CreateUser)
}

func Run(addr string) error {
    return r.Run(addr)
}
