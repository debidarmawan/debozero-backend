package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthenticatedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			c.Abort()
			return
		}

		tokenSplit := strings.Split(token, " ")
		if len(tokenSplit) != 2 || strings.ToLower(tokenSplit[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token, bearer token expected"})
			c.Abort()
			return
		}

		userID, err := tokenController.VerifyToken(tokenSplit[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
	}
}
