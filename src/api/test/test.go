package test

import (
	"net/http"

	t "Yosyos/src/model/const"

	"github.com/gin-gonic/gin"
)

func PassAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		password := c.GetHeader("Authorization")

		if password != t.Expected {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			c.Abort()
			return
		}

		c.Next()
	}
}
