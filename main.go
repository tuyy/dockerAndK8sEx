package main

import "os"
import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong22",
		})
	})

	r.GET("/hostname", func(c *gin.Context) {
        hostname, _ := os.Hostname()
		c.JSON(200, gin.H{
			"hostname": hostname,
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
