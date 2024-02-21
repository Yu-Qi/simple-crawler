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
	r.LoadHTMLGlob("templates/*")

	r.Use(checkHeadersMiddleware())

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main Page",
		})
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
