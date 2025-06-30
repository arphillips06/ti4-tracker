package main

import (
	"github.com/arphillips06/TI4-stats/database"
	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDatabase()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "TI4 stats tracker running"})
	})

	r.Run(":8080")
}
