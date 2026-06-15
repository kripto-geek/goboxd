package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Language string     `json:"language"`
	Source   string     `json:"source"`
	Tests    []TestPart `json:"tests"`
}

type TestPart struct {
	Stdin    string `json:"stdin"`
	Expected string `json:"expected"`
}

var sem = make(chan struct{}, 5)

func main() {
	r := gin.Default()

	r.POST("/run", RunHandler)
	r.Run(":8080")
}

func RunPythonCode(req Request, c *gin.Context) {
	dir, err := os.MkdirTemp("", "goboxd-*")
	if err != nil {
		return
	}

	defer os.RemoveAll(dir)

	sourcePath := filepath.Join(dir, "sol.py")

	os.WriteFile(sourcePath, []byte(req.Source), 0644)

	for _, test := range req.Tests {

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		cmd := exec.CommandContext(ctx, "python3", "sol.py")

		cmd.Stdin = strings.NewReader(test.Stdin)

		cmd.Dir = dir

		output, err := cmd.CombinedOutput()
		if err != nil {
			cancel()
			return
		}
		cancel()
		ActualOutput := strings.TrimSpace(string(output))
		ExpectedOutput := strings.TrimSpace(string(test.Expected))

		fmt.Print(ActualOutput)
		fmt.Print(ExpectedOutput)

		if ActualOutput != ExpectedOutput {
			c.JSON(200, gin.H{
				"output":   "Wrong Answer",
				"expected": ExpectedOutput,
				"acutal":   ActualOutput,
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"result": "Tests Success!",
	})
}

func RunCCode(req Request, c *gin.Context) {
	dir, err := os.MkdirTemp("", "goboxd-*")
	if err != nil {
		return
	}

	defer os.RemoveAll(dir)

	sourcePath := filepath.Join(dir, "sol.c")

	os.WriteFile(sourcePath, []byte(req.Source), 0644)

	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)

	cmd := exec.CommandContext(ctx, "gcc", "-o", "sol", "sol.c")

	cmd.Dir = dir

	cmd.CombinedOutput()

	for _, test := range req.Tests {

		cmd := exec.CommandContext(ctx, ".\\sol.exe")

		cmd.Stdin = strings.NewReader(test.Stdin)

		cmd.Dir = dir

		output, err := cmd.CombinedOutput()
		if err != nil {
			cancel()
			return
		}
		cancel()
		ActualOutput := strings.TrimSpace(string(output))
		ExpectedOutput := strings.TrimSpace(string(test.Expected))

		fmt.Print(ActualOutput)
		fmt.Print(ExpectedOutput)

		if ActualOutput != ExpectedOutput {
			c.JSON(200, gin.H{
				"output":   "Wrong Answer",
				"expected": ExpectedOutput,
				"acutal":   ActualOutput,
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"result": "Tests Success!",
	})
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
	fmt.Println(req.Language)

	sem <- struct{}{}
	defer func() { <-sem }()

	switch req.Language {

	case "python":
		RunPythonCode(req, c)

	case "C":
		RunCCode(req, c)

	}
}
