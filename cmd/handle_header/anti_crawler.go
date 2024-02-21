package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// checkHeadersMiddleware .
func checkHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAgent := c.GetHeader("User-Agent")
		referer := c.GetHeader("Referer")

		// verify the header
		if userAgent == "" || referer == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Access denied"})
			return
		}

		c.Next()
	}
}
