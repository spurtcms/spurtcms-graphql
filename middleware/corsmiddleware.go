package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS middleware
func CorsMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {

			c.AbortWithStatus(http.StatusNoContent)

			return
		}

		c.Next()
	}
}