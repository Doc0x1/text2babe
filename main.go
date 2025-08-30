package main

import (
	"doc0x1/text2babe/cmd"

	"github.com/joho/godotenv"
)

func init() {
	// Load .env file before any package initialization
	if err := godotenv.Load(); err != nil {
		// Silently ignore if no .env file - this is optional
	}
}

func main() {
	cmd.Execute()
}
