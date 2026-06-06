package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Source string `json:"source"`
}

func main() {
	r := gin.Default()

	r.POST("/run", RunHandler)
	r.Run(":8080")
}

func RunHandler(c *gin.Context) {
	var req Request
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Println(req.Source)
}
