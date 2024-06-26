package user

import (
    "net/http"
    "github.com/gin-gonic/gin"
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
    c.SetCookie("jwt", u.accessToken, 3600, "/", "localhost", false, true)
    c.JSON(http.StatusOK, u)
}

func (h *Handler) Logout(c *gin.Context) {
    // Reset cookie
    c.SetCookie("jwt", "", -1, "", "", false, true)
    c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
