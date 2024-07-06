package user

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/ihgazi/go-chat/config"
)

// Handler object is used to create the user 
// creation endpoint which is passed to GIN

type Handler struct {
    Service
}

func NewHandler(s Service) *Handler {
    return &Handler{
        Service: s,
    }
}

func (h *Handler) CreateUser(c *gin.Context) {
    var u CreateUserReq
    if err := c.ShouldBindJSON(&u); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Calling Service method
    res, err := h.Service.CreateUser(c.Request.Context(), &u)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, res)
}

func (h *Handler) Login(c *gin.Context) {
    var user LoginUserReq
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Calling service method
    u, err := h.Service.Login(c.Request.Context(), &user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Set cookie in context with JWT token
    domain := config.LoadConfig().ClientDomain
    secure := domain != "localhost"
    c.SetCookie("jwt", u.accessToken, 3600, "/", domain, secure, true)
    c.JSON(http.StatusOK, u)
}

func (h *Handler) Logout(c *gin.Context) {
    // Reset cookie
    domain := config.LoadConfig().ClientDomain
    secure := domain != "localhost"
    c.SetCookie("jwt", "", -1, "", domain, secure, true)
    c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Authentication endpoint to check validity of jwt token
// Can be reused for refresh token in future
func (h *Handler) AuthUser(c *gin.Context) {
    userID, _ := c.Get("userID")
    username, _ := c.Get("username")

    c.JSON(http.StatusOK, gin.H{"id": userID, "username": username})
}
