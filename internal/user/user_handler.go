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
