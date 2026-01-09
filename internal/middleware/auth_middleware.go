package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"dothefortune_server/internal/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		cookieToken, err := c.Cookie("token")
		if err == nil && cookieToken != "" {
			tokenString = cookieToken
		} else {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
				if tokenString == authHeader {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
					c.Abort()
					return
				}
			}
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}

