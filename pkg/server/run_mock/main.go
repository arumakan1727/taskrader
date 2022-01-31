package main

import (
	"github.com/arumakan1727/taskrader/pkg/server"
	"github.com/gin-gonic/gin"
)

func main() {
	r := server.NewEngine()

	r.GET("/taskrader", func(c *gin.Context) {
		c.File("../../../assets/index.html")
	})
	r.GET("/file/main.js", func(c *gin.Context) {
		c.File("../../../assets/main.js")
	})

	host := "localhost"
	port := ":8777"
	r.Run(host + port)
}
