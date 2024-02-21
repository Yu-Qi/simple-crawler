package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	port = 8555
)

func main() {
	r := gin.Default()
	r.Use(rateLimiter)

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome!")
	})

	go func() {
		// util server is running
		for {
			time.Sleep(1 * time.Second)
			_, err := http.Get(fmt.Sprintf("http://localhost:%d/", port))
			if err != nil {
				fmt.Println("Server check failed:", err)
				continue
			}
			break
		}
		crawler()
	}()

	r.Run(fmt.Sprintf(":%d", port))
}
