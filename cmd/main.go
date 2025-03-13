package main

import (
	app "crypto-trading-bot/internal"
	"log"
)

func main() {
	a := app.NewApp()
	if err := a.Run(); err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}
}
