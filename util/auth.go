package util

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ihgazi/go-chat/config"
)

// Custom JWT Claims class
type MyJWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func deauthorize(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	c.Abort()
}

// Auth middle to check validity of JWT token
func JWTValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := config.LoadConfig().SecretKey
		// Extract the token from the cookie
		tokenString, err := c.Cookie("jwt")
		if err != nil {
			deauthorize(c)
			return
		}

		// Parse and verify the token
		token, err := jwt.ParseWithClaims(tokenString, &MyJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil {
			deauthorize(c)
			return
		}

		if claims, ok := token.Claims.(*MyJWTClaims); ok && token.Valid {
			// Token is valid, store the claims in the context
			c.Set("userID", claims.ID)
			c.Set("username", claims.Username)
		} else {
			deauthorize(c)
			return
		}

		c.Next()
	}
}

// Create JWT token
func JWTCreateToken(ID string, Username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		ID:       ID,
		Username: Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // expires in one day
		},
	})

	// Fetch secret key from environment
	// sign token with key
	secretKey := config.LoadConfig().SecretKey
	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return ss, nil
}
