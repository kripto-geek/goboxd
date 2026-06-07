package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Source string `json:"source"`
	Stdin  string `json:"stdin"`
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

	dir, err := os.MkdirTemp("", "goboxd-*")
	if err != nil {
		return
	}

	defer os.RemoveAll(dir)

	sourcePath := filepath.Join(dir, "sol.py")

	os.WriteFile(sourcePath, []byte(req.Source), 0644)

	exec := exec.Command("python3", "sol.py")

	exec.Stdin = strings.NewReader(req.Stdin)

	exec.Dir = dir

	output, err := exec.CombinedOutput()
	if err != nil {
		return
	}

	c.JSON(200, gin.H{
		"output": string(output),
	})
}
