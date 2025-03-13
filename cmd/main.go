package main

import (
	"crypto-trading-bot/internal/app"
	"log"
)

func main() {
	a := app.NewApp()
	if err := a.Run(); err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}
}
